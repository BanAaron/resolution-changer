// this code was written on Christmas day because I have no life 😀

package main

import (
	"fmt"
	"github.com/banaaron/resolution-changer/displayManager"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"log/slog"
	"os"
)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func getIcon(fileLocation string) []byte {
	slog.Info("getting icon")
	iconBytes, err := os.ReadFile(fileLocation)
	if err != nil {
		slog.Error("failed to load icon", "error:", err)
	}
	return iconBytes
}

func onReady() {
	slog.Info("onReady")

	appName := "Resolution Changer"
	icon := getIcon("assets/icon_ico.ico")
	systray.SetIcon(icon)
	systray.SetTitle(appName)
	systray.SetTooltip(appName)

	// resolutions
	_3840x1080 := systray.AddMenuItem("3840x1080 (32:9)", "3840x1080")
	_2560x1080 := systray.AddMenuItem("2560x1080 (21:9)", "2560x1080")
	_1920x1080 := systray.AddMenuItem("1920x1080 (16:9)", "1920x1080")
	// refresh rates
	refreshRate := systray.AddMenuItem("Refresh Rate", "refresh rate")
	_144 := refreshRate.AddSubMenuItem("144hz", "144")
	_120 := refreshRate.AddSubMenuItem("120hz", "120")
	_75 := refreshRate.AddSubMenuItem("75hz", "75")
	_60 := refreshRate.AddSubMenuItem("60hz", "60")
	// separator
	systray.AddSeparator()
	// exit
	quit := systray.AddMenuItem("Exit", "exit")

	// create a goroutine
	go func() {
		var err error
		// infinite loop
		for {
			// select listens for all channels
			select {
			case <-_3840x1080.ClickedCh:
				err = displayManager.ChangeResolution(3840, 1080)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_2560x1080.ClickedCh:
				err = displayManager.ChangeResolution(2560, 1080)

				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_1920x1080.ClickedCh:
				err = displayManager.ChangeResolution(1920, 1080)
				panicError(err)

				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_144.ClickedCh:
				err = displayManager.ChangeRefreshRate(144)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_120.ClickedCh:
				err = displayManager.ChangeRefreshRate(120)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_75.ClickedCh:
				err = displayManager.ChangeRefreshRate(75)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-_60.ClickedCh:
				err = displayManager.ChangeRefreshRate(60)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					err = beeep.Notify("Error", errorString, "assets/icon_ico.ico")
					panicError(err)
				}
			case <-quit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	slog.Info("onExit")
}

func main() {
	systray.Run(onReady, onExit)
}
