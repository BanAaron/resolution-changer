// this code was written on Christmas day because I have no life 😀

package main

import (
	"fmt"
	"github.com/banaaron/resolution-changer/displayManager"
	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"log/slog"
	"os"
	"sync"
	"time"
)

type resMenu struct {
	item *systray.MenuItem
	res  displayManager.Resolution
}

type rateMenu struct {
	item *systray.MenuItem
	rate displayManager.RefreshRate
}

var (
	resMenus  []resMenu
	rateMenus []rateMenu

	currentRes  displayManager.Resolution
	currentRate displayManager.RefreshRate

	stateMu sync.Mutex
)

func applyDisplayInfo(di displayManager.DisplayInfo) {
	stateMu.Lock()
	defer stateMu.Unlock()

	for _, rm := range resMenus {
		if rm.res.Width == di.Resolution.Width && rm.res.Height == di.Resolution.Height {
			rm.item.Check()
		} else {
			rm.item.Uncheck()
		}
	}

	for _, hm := range rateMenus {
		if hm.rate == di.Refresh {
			hm.item.Check()
		} else {
			hm.item.Uncheck()
		}
	}

	currentRes = di.Resolution
	currentRate = di.Refresh
}

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

	// build resolution menu
	resMenus = nil
	for _, r := range cfg.Resolutions {
		label := fmt.Sprintf("%dx%d", r.Width, r.Height)
		item := systray.AddMenuItem(label, label)
		resMenus = append(resMenus, resMenu{
			item: item,
			res:  r,
		})
	}

	systray.AddSeparator()

	// build refresh rate menu
	rateMenus = nil
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

	// initial state: mark current res / Hz
	if di, err := displayManager.GetCurrentDisplay(); err != nil {
		slog.Error("GetCurrentDisplay failed", "err", err)
	} else {
		applyDisplayInfo(di)
	}

	// resolution handlers
	for _, rm := range resMenus {
		rm := rm
		go func() {
			for range rm.item.ClickedCh {
				err := displayManager.ChangeResolution(rm.res)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					if nErr := beeep.Notify("Error", errorString, iconLocation); nErr != nil {
						panicError(nErr)
					}
					continue
				}

				if di, e := displayManager.GetCurrentDisplay(); e == nil {
					applyDisplayInfo(di)
				}
			}
		}()
	}

	// refresh rate handlers
	for _, hm := range rateMenus {
		hm := hm
		go func() {
			for range hm.item.ClickedCh {
				err := displayManager.ChangeRefreshRate(hm.rate)
				if err != nil {
					errorString := fmt.Sprintf("%v", err)
					if nErr := beeep.Notify("Error", errorString, iconLocation); nErr != nil {
						panicError(nErr)
					}
					continue
				}

				if di, e := displayManager.GetCurrentDisplay(); e == nil {
					applyDisplayInfo(di)
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

	// "listener": poll for external changes every few seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)

			di, err := displayManager.GetCurrentDisplay()
			if err != nil {
				continue
			}

			stateMu.Lock()
			same := di.Resolution.Width == currentRes.Width &&
				di.Resolution.Height == currentRes.Height &&
				di.Refresh == currentRate
			stateMu.Unlock()

			if !same {
				applyDisplayInfo(di)
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
