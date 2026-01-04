package hidutils

import "testing"

func TestParseChargingState(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected ChargingState
	}{
		{"Charging state", 1, ChargingStateCharging},
		{"Discharging state", 3, ChargingStateDischarging},
		{"Disconnected state (0)", 0, ChargingStateDisconnected},
		{"Disconnected state (2)", 2, ChargingStateDisconnected},
		{"Unknown value 4", 4, ChargingStateDisconnected},
		{"Unknown value 255", 255, ChargingStateDisconnected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseChargingState(tt.input)
			if result != tt.expected {
				t.Errorf("ParseChargingState(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestChargingStateString(t *testing.T) {
	tests := []struct {
		state    ChargingState
		expected string
	}{
		{ChargingStateCharging, "Charging"},
		{ChargingStateDischarging, "Discharging"},
		{ChargingStateDisconnected, "Disconnected"},
		{ChargingStateUnknown, "Unknown"},
		{ChargingState(999), "Unknown"}, // Invalid value
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("ChargingState(%d).String() = %q, want %q", tt.state, result, tt.expected)
			}
		})
	}
}
