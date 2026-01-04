package hidutils

import "testing"

func TestParseBatteryLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected BatteryLevel
	}{
		{"Full battery", 4, BatteryLevelFull},
		{"Medium battery", 3, BatteryLevelMiddle},
		{"Low battery", 2, BatteryLevelLow},
		{"Empty battery", 1, BatteryLevelEmpty},
		{"Unknown value 0", 0, BatteryLevelEmpty},
		{"Unknown value 5", 5, BatteryLevelEmpty},
		{"Unknown value 255", 255, BatteryLevelEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBatteryLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseBatteryLevel(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBatteryLevelString(t *testing.T) {
	tests := []struct {
		level    BatteryLevel
		expected string
	}{
		{BatteryLevelFull, "Full"},
		{BatteryLevelMiddle, "Good"},
		{BatteryLevelLow, "Low"},
		{BatteryLevelEmpty, "Empty"},
		{BatteryLevelUnknown, "Unknown"},
		{BatteryLevel(999), "Unknown"}, // Invalid value
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("BatteryLevel(%d).String() = %q, want %q", tt.level, result, tt.expected)
			}
		})
	}
}
