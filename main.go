package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

var appName = "Resolution Changer"

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle(appName)
	systray.SetTooltip(appName)
}
func onExit() {}

func main() {
	systray.Run(onReady, onExit)
}
