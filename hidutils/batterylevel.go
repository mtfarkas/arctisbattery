package hidutils

// BatteryTreshold represents a battery threshold range.
// Not used yet, but here if we need it.
type BatteryTreshold struct {
	Min int
	Max int
}

// BatteryLevel is how charged the battery is.
type BatteryLevel int

const (
	// BatteryLevelUnknown means we couldn't read the level.
	BatteryLevelUnknown BatteryLevel = iota
	// BatteryLevelEmpty means battery is dead or nearly dead.
	BatteryLevelEmpty
	// BatteryLevelLow means battery is getting low.
	BatteryLevelLow
	// BatteryLevelMiddle means battery is doing alright.
	BatteryLevelMiddle
	// BatteryLevelFull means battery is full or nearly full.
	BatteryLevelFull
)

// String returns the battery level as text.
func (b BatteryLevel) String() string {
	switch b {
	case BatteryLevelFull:
		return "Full"
	case BatteryLevelMiddle:
		return "Good"
	case BatteryLevelLow:
		return "Low"
	case BatteryLevelEmpty:
		return "Empty"
	default:
		return "Unknown"
	}
}

// ParseBatteryLevel converts the device's byte value to a BatteryLevel.
func ParseBatteryLevel(input byte) BatteryLevel {
	switch input {
	case 4:
		return BatteryLevelFull
	case 3:
		return BatteryLevelMiddle
	case 2:
		return BatteryLevelLow
	default:
		return BatteryLevelEmpty
	}
}
