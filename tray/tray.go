// Package tray handles the system tray icon and menu.
package tray

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
	"github.com/mtfarkas/arctisbattery/config"
	"github.com/mtfarkas/arctisbattery/hidutils"
	"github.com/mtfarkas/arctisbattery/icons"
	"github.com/mtfarkas/arctisbattery/logging"
)

var doneCtx, doneCanceller = context.WithCancel(context.Background())

// Done returns a channel that closes when the tray shuts down.
func Done() <-chan struct{} {
	return doneCtx.Done()
}

func onStart() {
	initializeDevices()

	initializeTrayMenu()

	go updateLoop()
}

func onStop() {
	for _, d := range availableDevices {
		d.Raw.Device.Close()
	}
	availableDevices = nil
	doneCanceller()
}

const maxPanicRetries = 3

func updateLoop() {
	panicCount := 0
	backoff := 1 * time.Second

	for {
		exited := runUpdateLoop()
		if exited {
			return
		}

		// runUpdateLoop only returns false on panic recovery
		panicCount++
		if panicCount >= maxPanicRetries {
			log.Printf("updateLoop: reached max panic retries (%d), giving up", maxPanicRetries)
			return
		}

		log.Printf("updateLoop: recovered from panic (%d/%d), retrying in %v", panicCount, maxPanicRetries, backoff)
		select {
		case <-time.After(backoff):
			backoff *= 2
		case <-Done():
			return
		}
	}
}

// runUpdateLoop runs the polling loop. Returns true for normal exit, false if it recovered from a panic.
func runUpdateLoop() (normalExit bool) {
	defer func() {
		if r := recover(); r != nil {
			logging.RecoverAndLogValue("updateLoop", r)
			normalExit = false
		}
	}()

	for {
		queryDevicesAndUpdateValues()
		timer := time.NewTimer(300 * time.Millisecond)

		select {
		case <-timer.C:
			timer.Stop()
			continue
		case <-Done():
			timer.Stop()
			return true
		}
	}
}

func queryDevicesAndUpdateValues() {
	cfg := config.Get()
	if cfg == nil {
		log.Println("Warning: config not available")
		return
	}
	notifyEnabled := cfg.EnableNotify

	for _, d := range availableDevices {
		status, level, err := hidutils.ReadDeviceBatteryInfo(d.Raw.Device)
		if err != nil {
			log.Printf("Error reading battery info for device %s: %v", d.Raw.Info.ProductStr, err)
			continue
		}

		// send notification if battery just went empty
		if notifyEnabled && level == hidutils.BatteryLevelEmpty && d.BatteryLevel != hidutils.BatteryLevelEmpty {
			if err := beeep.Alert("Arctis Battery Utility", fmt.Sprintf("Battery low: %s", d.Raw.Info.ProductStr), ""); err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}

		d.BatteryStatus = status
		d.BatteryLevel = level
		updateMenuEntry(d)
	}

	updateSystrayTooltipAndIcon()
}

func updateSystrayTooltipAndIcon() {
	iconFinal := false
	var iconKey icons.IconKey = icons.IconDisconnected
	lowestBatteryLevel := hidutils.BatteryLevelFull

	lines := make([]string, 0, len(availableDevices))

	for _, d := range availableDevices {
		if d.IsAvailable() {
			lines = append(lines, d.String())

			// show charging icon if any device is charging but not full
			if d.BatteryStatus == hidutils.ChargingStateCharging && d.BatteryLevel != hidutils.BatteryLevelFull {
				iconKey = icons.IconChargingBattery
				iconFinal = true
			}

			// find the lowest battery level
			if d.BatteryLevel <= lowestBatteryLevel && !iconFinal {
				switch d.BatteryLevel {
				case hidutils.BatteryLevelEmpty:
					iconKey = icons.IconEmptyBattery
				case hidutils.BatteryLevelFull:
					iconKey = icons.IconFullBattery
				case hidutils.BatteryLevelLow:
					iconKey = icons.IconLowBattery
				case hidutils.BatteryLevelMiddle:
					iconKey = icons.IconMediumBattery
				}
			}
		}
	}

	if len(lines) < 1 {
		lines = append(lines, "No devices connected.")
	}

	icon, err := icons.GetIcon(iconKey, true)
	if err != nil {
		log.Printf("Failed to load icon %s: %v", iconKey, err)
	} else {
		systray.SetIcon(*icon)
	}

	tooltipStr := fmt.Sprintf("Arctis Battery Utility\n%s", strings.Join(lines, "\n"))
	systray.SetTooltip(tooltipStr)
}

// Start launches the system tray in the background.
func Start() {
	go func() {
		defer logging.RecoverAndLog("systray.Run", nil)
		systray.Run(onStart, onStop)
	}()
}
