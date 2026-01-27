package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/banaaron/resolution-changer/displayManager"
)

type AppConfig struct {
	Resolutions  []displayManager.Resolution
	RefreshRates []displayManager.RefreshRate
}

func loadConfig(path string) (AppConfig, error) {
	var cfg AppConfig

	f, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		switch currentSection {
		case "Resolutions":
			parts := strings.Split(line, "x")
			if len(parts) != 2 {
				slog.Warn("invalid resolution line", "line", line)
				continue
			}
			w, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
			h, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err1 != nil || err2 != nil {
				slog.Warn("invalid resolution values", "line", line)
				continue
			}
			cfg.Resolutions = append(cfg.Resolutions, displayManager.Resolution{
				Width:  uint32(w),
				Height: uint32(h),
			})

		case "RefreshRates":
			hz, err := strconv.Atoi(line)
			if err != nil {
				slog.Warn("invalid refresh rate", "line", line)
				continue
			}
			cfg.RefreshRates = append(cfg.RefreshRates, displayManager.RefreshRate(hz))
		default:
			// ignore other sections
		}
	}

	if err := scanner.Err(); err != nil {
		return cfg, fmt.Errorf("scan config: %w", err)
	}

	return cfg, nil
}
