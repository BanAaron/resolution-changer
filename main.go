// this code was written on Christmas day because I have no life 😀

package main

import (
	"github.com/getlantern/systray"
	"log/slog"
	"os"
)

func getIcon(fileLocation string) []byte {
	b, err := os.ReadFile(fileLocation)
	if err != nil {
		slog.Error("%v", err)
	}
	return b
}

func changeResolution(width int, height int) {
	slog.Info("changeResolution", "width", width, "height", height)
}

func changeRefreshRate(refreshRate int) {
	slog.Info("changeRefreshRate", "refreshRate", refreshRate)
}

func onReady() {
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
	// exit
	quit := systray.AddMenuItem("Exit", "exit")

	// create a goroutine
	go func() {
		// infinite loop
		for {
			// select listens for all channels
			select {
			case <-_3840x1080.ClickedCh:
				changeResolution(3840, 1080)
			case <-_2560x1080.ClickedCh:
				changeResolution(2560, 1080)
			case <-_1920x1080.ClickedCh:
				changeResolution(1920, 1080)
			case <-_144.ClickedCh:
				changeRefreshRate(144)
			case <-_120.ClickedCh:
				changeRefreshRate(120)
			case <-_75.ClickedCh:
				changeRefreshRate(75)
			case <-_60.ClickedCh:
				changeRefreshRate(60)
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
