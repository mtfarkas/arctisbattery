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

func initializeDevices() {
	rawDevices := make([]*rawDevice, 0)

	err := hid.Enumerate(hidutils.SteelSeriesVID, hid.ProductIDAny, func(info *hid.DeviceInfo) error {
		device, err := hid.OpenPath(info.Path)

		// skip devices we can't open
		if err != nil {
			log.Printf("Could not open device %s: %v", info.ProductStr, err)
			return nil
		}

		// skip devices that don't support battery info
		_, _, err = hidutils.ReadDeviceBatteryInfo(device)
		if err != nil {
			log.Printf("Device %s (%s) does not support battery info: %v", info.ProductStr, info.Path, err)
			device.Close()
			return nil
		}

		log.Printf("Found compatible device: %s", info.ProductStr)
		rawDevices = append(rawDevices, &rawDevice{
			Device: device,
			Info:   info,
		})

		return nil
	})

	if err != nil {
		log.Printf("Error enumerating HID devices: %v", err)
	}

	for _, d := range rawDevices {
		availableDevices = append(availableDevices, createDevice(d))
	}

	log.Printf("Initialized %d compatible device(s)", len(availableDevices))
}

func createDevice(rawDevice *rawDevice) *deviceEntry {
	entry := &deviceEntry{
		Raw:           rawDevice,
		BatteryLevel:  hidutils.BatteryLevelUnknown,
		BatteryStatus: hidutils.ChargingStateUnknown,
		MenuItem:      nil,
	}
	return entry
}
