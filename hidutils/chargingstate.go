package hidutils

// ChargingState is whether the device is charging or not.
type ChargingState int

const (
	// ChargingStateUnknown means we couldn't read the state.
	ChargingStateUnknown ChargingState = iota
	// ChargingStateDisconnected means device isn't connected.
	ChargingStateDisconnected
	// ChargingStateCharging means device is plugged in and charging.
	ChargingStateCharging
	// ChargingStateDischarging means device is running on battery.
	ChargingStateDischarging
)

// String returns the charging state as text.
func (c ChargingState) String() string {
	switch c {
	case ChargingStateCharging:
		return "Charging"
	case ChargingStateDischarging:
		return "Discharging"
	case ChargingStateDisconnected:
		return "Disconnected"
	default:
		return "Unknown"
	}
}

// ParseChargingState converts the device's byte value to a ChargingState.
func ParseChargingState(input byte) ChargingState {
	switch input {
	case 1:
		return ChargingStateCharging
	case 3:
		return ChargingStateDischarging
	default:
		return ChargingStateDisconnected
	}
}
