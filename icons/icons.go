package icons

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

//go:embed assets/*
var embeddedAssets embed.FS

type IconKey string

const (
	IconEmptyBattery    IconKey = "assets/battery"
	IconLowBattery      IconKey = "assets/battery-low"
	IconMediumBattery   IconKey = "assets/battery-medium"
	IconFullBattery     IconKey = "assets/battery-full"
	IconChargingBattery IconKey = "assets/battery-charging"
	IconDisconnected    IconKey = "assets/unplug"
)

var iconsCache map[string][]byte = make(map[string][]byte)

func Init() {
	entries, err := embeddedAssets.ReadDir("assets")
	if err != nil {
		return
	}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != platformExtension() {
			continue
		}
		iconBytes, _ := embeddedAssets.ReadFile(path.Join("assets", entry.Name()))
		cacheKey := fmt.Sprintf("assets/%s", strings.Split(entry.Name(), ".")[0])
		iconsCache[cacheKey] = iconBytes
	}
}

func GetIcon(icon IconKey, darkMode bool) (*[]byte, error) {
	suffix := "dark"
	if !darkMode {
		suffix = "light"
	}
	cacheKey := fmt.Sprintf("%s-%s", icon, suffix)

	if val, ok := iconsCache[cacheKey]; ok {
		return &val, nil
	}

	return nil, fmt.Errorf("cache key %s not found", cacheKey)
}
