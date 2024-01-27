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

// refreshRate is the refresh rate in hz per second
type refreshRate uint32

var (
	user32dll                  = syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW   = user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW = user32dll.NewProc("ChangeDisplaySettingsW")
	devMode                    = new(devmode)
)

func ChangeResolution(res Resolution) error {
	var err error
	slog.Info("changeResolution", "Width", res.Width, "Height", res.Height)
	// get the display information
	response, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)), uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(devMode)))
	if response == 0 {
		err = fmt.Errorf("could not extract display settings")
		return err
	}

	// change the display Resolution
	newMode := *devMode
	newMode.DmPelsWidth = res.Width
	newMode.DmPelsHeight = res.Height
	response, _, _ = procChangeDisplaySettingsW.Call(uintptr(unsafe.Pointer(&newMode)), uintptr(0))

	switch response {
	case uintptr(DISP_CHANGE_SUCCESSFUL):
		slog.Info("successfully changed the display Resolution")
	case uintptr(DISP_CHANGE_RESTART):
		slog.Info("restart required to apply the Resolution changes")
	case uintptr(DISP_CHANGE_BADMODE):
		slog.Error("the Resolution is not supported by the display")
	case uintptr(DISP_CHANGE_FAILED):
		slog.Error("failed to change the display Resolution")
	}

	return err
}

func ChangeRefreshRate(refreshRate refreshRate) error {
	var err error
	slog.Info("ChangeRefreshRate", "refreshRate", refreshRate)
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
