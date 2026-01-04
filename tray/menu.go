package tray

import (
	"log"

	"fyne.io/systray"
	"github.com/mtfarkas/arctisbattery/config"
	"github.com/mtfarkas/arctisbattery/icons"
)

func initializeTrayMenu() {
	systray.SetTitle("Arctis Battery Utility")
	icon, err := icons.GetIcon(icons.IconDisconnected, true)
	if err == nil {
		systray.SetIcon(*icon)
	}

	systray.AddMenuItem("Devices", "A list of connected devices").Disable()

	for _, v := range availableDevices {
		v.MenuItem = systray.AddMenuItem("", "")
		updateMenuEntry(v)
	}

	systray.AddSeparator()
	setupNotificationButton()

	systray.AddSeparator()
	setupQuitButton()
}

func setupQuitButton() {
	mQuit := systray.AddMenuItem("Exit", "Quits the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func setupNotificationButton() {
	mNotify := systray.AddMenuItem("Notify on low battery", "")

	cfg := config.Get()
	if cfg != nil && cfg.EnableNotify {
		mNotify.Check()
	}

	go func() {
		for {
			select {
			case <-mNotify.ClickedCh:
				cfg := config.Get()
				if cfg == nil {
					log.Println("Warning: config not available")
					continue
				}

				if cfg.EnableNotify {
					mNotify.Uncheck()
					if err := config.Modify(func(ac *config.AppConfiguration) { ac.EnableNotify = false }); err != nil {
						log.Printf("Failed to modify config: %v", err)
					}
				} else {
					mNotify.Check()
					if err := config.Modify(func(ac *config.AppConfiguration) { ac.EnableNotify = true }); err != nil {
						log.Printf("Failed to modify config: %v", err)
					}
				}
			case <-Done():
				return
			}
		}
	}()
}

func updateMenuEntry(d *deviceEntry) {
	if !d.IsAvailable() {
		d.MenuItem.Hide()
	} else {
		d.MenuItem.SetTitle(d.String())
		d.MenuItem.Show()
	}
}
