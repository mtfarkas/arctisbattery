package config

import (
	"encoding/json"
	"os"
	"sync"
	"testing"
)

// resetConfig is a test helper that resets the config state
// Note: In production, you'd want to use dependency injection for better testability
func resetConfig() {
	configMutex.Lock()
	defer configMutex.Unlock()
	configInstance = nil
	configInitOnce = sync.Once{}
	configInitErr = nil
}

func TestConfigInit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	resetConfig()

	err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	cfg := Get()
	if cfg == nil {
		t.Fatal("Get() returned nil after successful Init()")
	}

	// Default value should be true
	if !cfg.EnableNotify {
		t.Error("Expected EnableNotify to be true by default")
	}
}

func TestConfigModify(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	resetConfig()

	err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Modify config
	err = Modify(func(ac *AppConfiguration) {
		ac.EnableNotify = false
	})
	if err != nil {
		t.Fatalf("Modify() failed: %v", err)
	}

	cfg := Get()
	if cfg.EnableNotify {
		t.Error("Expected EnableNotify to be false after modification")
	}
}

func TestConfigPersistence(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	resetConfig()

	// Initialize and modify
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	if err := Modify(func(ac *AppConfiguration) {
		ac.EnableNotify = false
	}); err != nil {
		t.Fatalf("Modify() failed: %v", err)
	}

	// Explicitly save for testing
	if err := save(); err != nil {
		t.Fatalf("save() failed: %v", err)
	}

	// Reset and reload
	resetConfig()

	if err := Init(); err != nil {
		t.Fatalf("Second Init() failed: %v", err)
	}

	cfg := Get()
	if cfg.EnableNotify {
		t.Error("Config not persisted correctly, expected EnableNotify to be false")
	}
}

func TestEnsureConfigExists(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// File shouldn't exist yet
	if _, err := os.Stat(configFilePath); err == nil {
		t.Fatal("Config file already exists before ensureConfigExists()")
	}

	err := ensureConfigExists()
	if err != nil {
		t.Fatalf("ensureConfigExists() failed: %v", err)
	}

	// File should now exist
	if _, err := os.Stat(configFilePath); err != nil {
		t.Fatal("Config file doesn't exist after ensureConfigExists()")
	}

	// Verify content
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if string(content) != "{}" {
		t.Errorf("Expected empty JSON object, got: %s", content)
	}
}

func TestLoadConfigWithInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create invalid JSON file
	invalidJSON := []byte("{invalid json}")
	if err := os.WriteFile(configFilePath, invalidJSON, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	resetConfig()

	err := Init()
	if err == nil {
		t.Error("Expected Init() to fail with invalid JSON")
	}
}

func TestConfigMarshaling(t *testing.T) {
	cfg := &AppConfiguration{
		EnableNotify: true,
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	var unmarshaled AppConfiguration
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	if unmarshaled.EnableNotify != cfg.EnableNotify {
		t.Error("Marshaling/unmarshaling changed EnableNotify value")
	}
}

func TestModifyWithoutInit(t *testing.T) {
	resetConfig()

	err := Modify(func(ac *AppConfiguration) {
		ac.EnableNotify = false
	})

	if err == nil {
		t.Error("Expected Modify() to fail when config is not initialized")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	resetConfig()

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Test concurrent reads and writes
	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				_ = Get()
			}
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		go func(val bool) {
			for j := 0; j < 10; j++ {
				Modify(func(ac *AppConfiguration) {
					ac.EnableNotify = val
				})
			}
			done <- true
		}(i%2 == 0)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic or deadlock
	cfg := Get()
	if cfg == nil {
		t.Error("Config became nil during concurrent access")
	}
}
