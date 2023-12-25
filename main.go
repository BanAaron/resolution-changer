package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"log/slog"
)

func onReady() {
	appName := "Resolution Changer"
	systray.SetIcon(icon.Data)
	systray.SetTitle(appName)
	systray.SetTooltip(appName)

	// menu items
	res := systray.AddMenuItem("Res", "change res")
	quit := systray.AddMenuItem("Exit", "exit the program")

	go func() {
		for {
			// select listens for all channels
			select {
			// when the quit button channel is updated we do something
			case <-quit.ClickedCh:
				systray.Quit()
			case <-res.ClickedCh:
				fmt.Println("res")
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
