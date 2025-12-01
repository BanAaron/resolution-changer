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

var (
	user32dll                  = syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW   = user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW = user32dll.NewProc("ChangeDisplaySettingsW")
	devMode                    = new(devmode)
)

func ChangeResolution(res Resolution) error {
	var err error
	slog.Info("changeResolution", "Width", res.Width, "Height", res.Height)

	var mode devmode
	mode.DmSize = uint16(unsafe.Sizeof(mode))

	// get available resolutions
	for i := 0; ; i++ {
		response, _, _ := procEnumDisplaySettingsW.Call(
			uintptr(unsafe.Pointer(nil)),
			uintptr(i),
			uintptr(unsafe.Pointer(&mode)),
		)
		if response == 0 {
			break
		}

		if mode.DmPelsWidth == res.Width && mode.DmPelsHeight == res.Height {
			// found a matching mode; try to set it
			response, _, _ = procChangeDisplaySettingsW.Call(
				uintptr(unsafe.Pointer(&mode)),
				uintptr(0),
			)

			switch response {
			case uintptr(DISP_CHANGE_SUCCESSFUL):
				slog.Info("successfully changed the display Resolution")
			case uintptr(DISP_CHANGE_RESTART):
				slog.Info("restart required to apply the Resolution changes")
				err = fmt.Errorf("restart required to apply resolution %dx%d", res.Width, res.Height)
			case uintptr(DISP_CHANGE_BADMODE):
				slog.Error("the Resolution is not supported by the display")
				err = fmt.Errorf("the Resolution %dx%d is not supported by the display (BADMODE)", res.Width, res.Height)
			case uintptr(DISP_CHANGE_FAILED):
				slog.Error("failed to change the display Resolution")
				err = fmt.Errorf("failed to change the display Resolution to %dx%d (FAILED)", res.Width, res.Height)
			default:
				err = fmt.Errorf("ChangeDisplaySettingsW returned unexpected code: %d", response)
			}

			return err
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
