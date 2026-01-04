package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mtfarkas/arctisbattery/config"
	"github.com/mtfarkas/arctisbattery/hidutils"
	"github.com/mtfarkas/arctisbattery/icons"
	"github.com/mtfarkas/arctisbattery/tray"
)

func main() {
	// load config first
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// set up HID library
	err := hidutils.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer hidutils.Exit()

	// load icons
	icons.Init()

	// start the tray icon
	tray.Start()

	// wait for shutdown signal
	sigIntCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-tray.Done():
		os.Exit(0)
	case <-sigIntCtx.Done():
		os.Exit(0)
		stop()
	}
}
