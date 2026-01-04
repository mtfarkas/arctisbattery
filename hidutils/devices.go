// Package hidutils talks to SteelSeries HID devices to read battery status.
package hidutils

import (
	"time"

	"github.com/sstallion/go-hid"
)

// SteelSeriesVID is the USB vendor ID for SteelSeries.
const SteelSeriesVID = 0x1038

// Init sets up the HID library. Call this before anything else.
func Init() error {
	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		return err
	}

	return nil
}

// Exit cleans up the HID library.
func Exit() error {
	return hid.Exit()
}

// ReadDeviceBatteryInfo reads battery info from a device.
func ReadDeviceBatteryInfo(d *hid.Device) (ChargingState, BatteryLevel, error) {
	batteryStatusCommand := []byte{0x00, 0xb0}

	if _, err := d.Write(batteryStatusCommand); err != nil {
		return ChargingStateUnknown, BatteryLevelUnknown, err
	}

	readBuf := make([]byte, 4)
	if _, err := d.ReadWithTimeout(readBuf, time.Millisecond*100); err != nil {
		return ChargingStateUnknown, BatteryLevelUnknown, err
	}

	chargingStatus := ParseChargingState(readBuf[3])
	batteryLevel := ParseBatteryLevel(readBuf[2])

	return chargingStatus, batteryLevel, nil
}
