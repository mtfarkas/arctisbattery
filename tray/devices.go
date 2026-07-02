package tray

import (
	"fmt"
	"log"

	"fyne.io/systray"
	"github.com/mtfarkas/arctisbattery/hidutils"
	"github.com/sstallion/go-hid"
)

type rawDevice struct {
	Device *hid.Device
	Info   *hid.DeviceInfo
}

type deviceEntry struct {
	Raw *rawDevice

	BatteryLevel  hidutils.BatteryLevel
	BatteryStatus hidutils.ChargingState

	MenuItem *systray.MenuItem
}

func (d *deviceEntry) IsAvailable() bool {
	return d.BatteryStatus != hidutils.ChargingStateDisconnected && d.BatteryStatus != hidutils.ChargingStateUnknown
}

func (d *deviceEntry) String() string {
	batteryStateStr := "↓"
	if d.BatteryStatus == hidutils.ChargingStateCharging {
		batteryStateStr = "↑"
	}

	return fmt.Sprintf("%s %s - %s", batteryStateStr, d.Raw.Info.ProductStr, d.BatteryLevel.String())
}

var availableDevices []*deviceEntry = make([]*deviceEntry, 0)

type initialDeviceState struct {
	raw    *rawDevice
	status hidutils.ChargingState
	level  hidutils.BatteryLevel
}

func initializeDevices() {
	var initialStates []initialDeviceState

	err := hid.Enumerate(hidutils.SteelSeriesVID, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		device, err := hid.OpenPath(info.Path)

		// skip devices we can't open
		if err != nil {
			log.Printf("Could not open device %s: %v", info.ProductStr, err)
			return nil
		}

		// skip devices that don't support battery info; capture the initial
		// reading so createDevice can seed the entry with a real prior state
		status, level, err := hidutils.ReadDeviceBatteryInfo(device)
		if err != nil {
			log.Printf("Device %s (%s) does not support battery info: %v", info.ProductStr, info.Path, err)
			device.Close()
			return nil
		}

		log.Printf("Found compatible device: %s", info.ProductStr)
		initialStates = append(initialStates, initialDeviceState{
			raw:    &rawDevice{Device: device, Info: info},
			status: status,
			level:  level,
		})

		return nil
	})

	if err != nil {
		log.Printf("Error enumerating HID devices: %v", err)
	}

	for _, s := range initialStates {
		availableDevices = append(availableDevices, createDevice(s.raw, s.status, s.level))
	}

	log.Printf("Initialized %d compatible device(s)", len(availableDevices))
}

func createDevice(rawDevice *rawDevice, initialStatus hidutils.ChargingState, initialLevel hidutils.BatteryLevel) *deviceEntry {
	return &deviceEntry{
		Raw:           rawDevice,
		BatteryLevel:  initialLevel,
		BatteryStatus: initialStatus,
		MenuItem:      nil,
	}
}
