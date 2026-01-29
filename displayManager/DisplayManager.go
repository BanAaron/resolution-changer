package displayManager

import (
	"fmt"
	"log/slog"
	"syscall"
	"unsafe"
)

// they are in shouting snake case to match the Windows api conventions
const (
	CCHDEVICENAME                 = 32
	CCHFORMNAME                   = 32
	ENUM_CURRENT_SETTINGS  uint32 = 0xFFFFFFFF
	DISP_CHANGE_SUCCESSFUL uint32 = 0
	DISP_CHANGE_RESTART    uint32 = 1
	DISP_CHANGE_FAILED     uint32 = 0xFFFFFFFF
	DISP_CHANGE_BADMODE    uint32 = 0xFFFFFFFE
	CDS_UPDATEREGISTRY     uint32 = 0x00000001
)

// devmode is a structure used to specify characteristics
// of display and print devices.
type devmode struct {
	DmDeviceName       [CCHDEVICENAME]uint16
	DmSpecVersion      uint16
	DmDriverVersion    uint16
	DmSize             uint16
	DmDriverExtra      uint16
	DmFields           uint32
	DmOrientation      int16
	DmPaperSize        int16
	DmPaperLength      int16
	DmPaperWidth       int16
	DmScale            int16
	DmCopies           int16
	DmDefaultSource    int16
	DmPrintQuality     int16
	DmColor            int16
	DmDuplex           int16
	DmYResolution      int16
	DmTTOption         int16
	DmCollate          int16
	DmFormName         [CCHFORMNAME]uint16
	DmLogPixels        uint16
	DmBitsPerPel       uint32
	DmPelsWidth        uint32
	DmPelsHeight       uint32
	DmDisplayFlags     uint32
	DmDisplayFrequency uint32
	DmICMMethod        uint32
	DmICMIntent        uint32
	DmMediaType        uint32
	DmDitherType       uint32
	DmReserved1        uint32
	DmReserved2        uint32
	DmPanningWidth     uint32
	DmPanningHeight    uint32
}

// Resolution is the Width and Height in pixels
type Resolution struct {
	Width, Height uint32
}

// RefreshRate is the refresh rate in hz per second
type RefreshRate uint32

type DisplayInfo struct {
	Resolution
	Refresh RefreshRate
}

var (
	user32dll                  = syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW   = user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW = user32dll.NewProc("ChangeDisplaySettingsW")
	devMode                    = new(devmode)
)

func GetCurrentDisplay() (DisplayInfo, error) {
	var dm devmode
	dm.DmSize = uint16(unsafe.Sizeof(dm))

	response, _, _ := procEnumDisplaySettingsW.Call(
		uintptr(unsafe.Pointer(nil)),
		uintptr(ENUM_CURRENT_SETTINGS),
		uintptr(unsafe.Pointer(&dm)),
	)
	if response == 0 {
		return DisplayInfo{}, fmt.Errorf("could not extract display settings")
	}

	return DisplayInfo{
		Resolution: Resolution{
			Width:  dm.DmPelsWidth,
			Height: dm.DmPelsHeight,
		},
		Refresh: RefreshRate(dm.DmDisplayFrequency),
	}, nil
}

func ChangeResolution(res Resolution) error {
	var err error
	slog.Info("changeResolution", "Width", res.Width, "Height", res.Height)

	// get current display settings (to read current refresh rate)
	var current devmode
	current.DmSize = uint16(unsafe.Sizeof(current))
	response, _, _ := procEnumDisplaySettingsW.Call(
		uintptr(unsafe.Pointer(nil)),
		uintptr(ENUM_CURRENT_SETTINGS),
		uintptr(unsafe.Pointer(&current)),
	)
	if response == 0 {
		return fmt.Errorf("could not extract display settings")
	}
	currentHz := current.DmDisplayFrequency

	// helper to apply a mode
	applyMode := func(m *devmode) error {
		response, _, _ := procChangeDisplaySettingsW.Call(
			uintptr(unsafe.Pointer(m)),
			uintptr(CDS_UPDATEREGISTRY),
		)

		switch response {
		case uintptr(DISP_CHANGE_SUCCESSFUL):
			slog.Info("successfully changed the display Resolution")
			return nil
		case uintptr(DISP_CHANGE_RESTART):
			slog.Info("restart required to apply the Resolution changes")
			return fmt.Errorf("restart required to apply resolution %dx%d", res.Width, res.Height)
		case uintptr(DISP_CHANGE_BADMODE):
			slog.Error("the Resolution is not supported by the display")
			return fmt.Errorf("the Resolution %dx%d is not supported by the display (BADMODE)", res.Width, res.Height)
		case uintptr(DISP_CHANGE_FAILED):
			slog.Error("failed to change the display Resolution")
			return fmt.Errorf("failed to change the display Resolution to %dx%d (FAILED)", res.Width, res.Height)
		default:
			return fmt.Errorf("ChangeDisplaySettingsW returned unexpected code: %d", response)
		}
	}

	var mode devmode
	mode.DmSize = uint16(unsafe.Sizeof(mode))

	// first try: match Resolution + RefreshRate
	for i := 0; ; i++ {
		response, _, _ := procEnumDisplaySettingsW.Call(
			uintptr(unsafe.Pointer(nil)),
			uintptr(i),
			uintptr(unsafe.Pointer(&mode)),
		)
		if response == 0 {
			break
		}

		if mode.DmPelsWidth == res.Width &&
			mode.DmPelsHeight == res.Height &&
			mode.DmDisplayFrequency == currentHz {
			return applyMode(&mode)
		}
	}

	// fallback: Resolution only if no Hz match
	mode.DmSize = uint16(unsafe.Sizeof(mode))
	for i := 0; ; i++ {
		response, _, _ := procEnumDisplaySettingsW.Call(
			uintptr(unsafe.Pointer(nil)),
			uintptr(i),
			uintptr(unsafe.Pointer(&mode)),
		)
		if response == 0 {
			break
		}

		if mode.DmPelsWidth == res.Width &&
			mode.DmPelsHeight == res.Height {
			return applyMode(&mode)
		}
	}

	err = fmt.Errorf("resolution %dx%d not found in EnumDisplaySettingsW", res.Width, res.Height)
	return err
}

func ChangeRefreshRate(refreshRate RefreshRate) error {
	var err error
	slog.Info("ChangeRefreshRate", "RefreshRate", refreshRate)
	// get the display information
	response, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)), uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(devMode)))
	if response == 0 {
		err = fmt.Errorf("could not extract display settings")
		return err
	}
	// change the display Resolution
	newMode := *devMode
	newMode.DmDisplayFrequency = uint32(refreshRate)
	response, _, _ = procChangeDisplaySettingsW.Call(uintptr(unsafe.Pointer(&newMode)), uintptr(0))

	switch response {
	case uintptr(DISP_CHANGE_SUCCESSFUL):
		slog.Info("successfully changed the display refresh rate")
	case uintptr(DISP_CHANGE_RESTART):
		slog.Info("restart required to apply the refresh rate changes")
	case uintptr(DISP_CHANGE_BADMODE):
		err = fmt.Errorf("the refresh rate is not supported by the display")
	case uintptr(DISP_CHANGE_FAILED):
		err = fmt.Errorf("failed to change the display refresh rate")
	}

	return err
}
