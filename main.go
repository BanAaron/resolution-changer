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
	iconLocation := "assets/icon_ico.ico"

	icon := getIcon(iconLocation)
	systray.SetIcon(icon)
	systray.SetTitle(appName)
	systray.SetTooltip(appName)

	cfg, err := loadConfig("config.ini")
	if err != nil {
		slog.Warn("using built-in defaults; failed to load config.ini", "err", err)
		cfg = AppConfig{
			Resolutions: []displayManager.Resolution{
				{Width: 2560, Height: 1600},
				{Width: 2560, Height: 1440},
			},
			RefreshRates: []displayManager.RefreshRate{
				240,
				60,
			},
		}
	}

	type resMenu struct {
		item *systray.MenuItem
		res  displayManager.Resolution
	}
	type rateMenu struct {
		item *systray.MenuItem
		rate displayManager.RefreshRate
	}

	var resMenus []resMenu
	for _, r := range cfg.Resolutions {
		label := fmt.Sprintf("%dx%d", r.Width, r.Height)
		item := systray.AddMenuItem(label, label)
		resMenus = append(resMenus, resMenu{
			item: item,
			res:  r,
		})
	}

	systray.AddSeparator()

	var rateMenus []rateMenu
	for _, hz := range cfg.RefreshRates {
		label := fmt.Sprintf("%dhz", hz)
		item := systray.AddMenuItem(label, label)
		rateMenus = append(rateMenus, rateMenu{
			item: item,
			rate: hz,
		})
	}

	systray.AddSeparator()
	quit := systray.AddMenuItem("Exit", "exit")

	// handlers for resolution items
	for _, rm := range resMenus {
		rm := rm // capture
		go func() {
			for range rm.item.ClickedCh {
				err := displayManager.ChangeResolution(rm.res)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					if nErr := beeep.Notify("Error", errorString, iconLocation); nErr != nil {
						panicError(nErr)
					}
				}
			}
		}()
	}

	// handlers for refresh rate items
	for _, hm := range rateMenus {
		hm := hm // capture
		go func() {
			for range hm.item.ClickedCh {
				err := displayManager.ChangeRefreshRate(hm.rate)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					if nErr := beeep.Notify("Error", errorString, iconLocation); nErr != nil {
						panicError(nErr)
					}
				}
			}
		}()
	}

	// quit handler
	go func() {
		for range quit.ClickedCh {
			systray.Quit()
		}
	}()
}

func onExit() {
	slog.Info("onExit")
}

func main() {
	systray.Run(onReady, onExit)
}
