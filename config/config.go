// Package config handles loading and saving settings.
package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

const configFilePath = "arctisbattery_settings.json"

// AppConfiguration is the settings struct.
type AppConfiguration struct {
	// EnableNotify toggles low battery notifications.
	EnableNotify bool `json:"enable_notify"`
}

var (
	configInstance *AppConfiguration
	configMutex    sync.RWMutex
	configInitOnce sync.Once
	configInitErr  error
)

// Init loads the config. Call this before Get() or Modify().
func Init() error {
	configInitOnce.Do(func() {
		configInitErr = loadConfig()
	})
	return configInitErr
}

// Get returns the current config. Returns nil if Init() wasn't called or failed.
func Get() *AppConfiguration {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return configInstance
}

type ConfigModifyFunc func(*AppConfiguration)

// Modify updates the config and saves it to disk asynchronously.
func Modify(f ConfigModifyFunc) error {
	configMutex.Lock()
	if configInstance == nil {
		configMutex.Unlock()
		return errors.New("config not initialized")
	}
	f(configInstance)
	configMutex.Unlock()

	// save to disk in background
	go func() {
		if err := save(); err != nil {
			log.Printf("Failed to save config: %v", err)
		}
	}()

	return nil
}

func loadConfig() error {
	if err := ensureConfigExists(); err != nil {
		return err
	}

	cfgBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	cfg := &AppConfiguration{
		EnableNotify: true,
	}
	if err := json.Unmarshal(cfgBytes, cfg); err != nil {
		return err
	}

	configInstance = cfg
	return nil
}

func save() error {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if configInstance == nil {
		return errors.New("config not initialized")
	}

	if err := ensureConfigExists(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(*configInstance, "", "  ")
	if err != nil {
		return err
	}

	// write with truncate + create if needed
	return os.WriteFile(configFilePath, data, 0644)
}

func ensureConfigExists() error {
	if _, err := os.Stat(configFilePath); err == nil || !errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("{}")
	return err
}
