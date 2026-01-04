# Arctis Battery Utility

Simple system tray app to check your SteelSeries headset battery level without installing the bloated GG software.

## What it does

Shows your SteelSeries headset battery in the Windows system tray. That's it. No ads, no account required, no telemetry.

- Real-time battery monitoring
- Low battery notifications (optional)
- System tray icon that changes based on charge level
- Lightweight (under 10MB RAM)
- Single executable, no installation

## Supported Headsets

Should work with most SteelSeries wireless headsets, depends on the protocol used.

Tested with:
- Arctis Nova 7X

## Installation

Grab the latest release from the releases page, or build it yourself:

```bash
git clone github.com/mtfarkas/arctisbattery
cd arctisbattery
make
```

You'll need Go 1.21+ and a C compiler (for the HID library).

## How to use

Just run it. The icon will appear in your system tray showing your current battery level. Click it to see connected devices and toggle notifications.

Settings are saved to `arctisbattery_settings.json` in the same folder as the executable.

## Development

Run tests with `go test ./...`. The code is organized into a few packages:
- `config` - handles settings
- `hidutils` - talks to the HID devices
- `icons` - icon loading
- `tray` - system tray UI

## How it works

I use the HID protocol SteelSeries uses. The app sends `0x00, 0xb0` to the device and reads back:
- Byte 3: Battery level (0-4)
- Byte 4: Charging state (1=charging, 3=discharging)

Pretty simple. No network calls, no telemetry, just reads from the USB device.

## Dependencies

- [go-hid](https://github.com/sstallion/go-hid) - for talking to USB HID devices
- [systray](https://github.com/fyne-io/systray) - for the tray icon
- [beeep](https://github.com/gen2brain/beeep) - for notifications

## Troubleshooting

**Device not showing up?**
- Make sure it's on and connected
- Try unplugging and replugging the USB dongle

## FAQ

**Does this void my warranty?**  
Nope. It only reads from the device, doesn't write anything or modify firmware.

**Why only 4 battery levels (Empty, Low, Good, Full)?**  
That's all the headset reports. Even the official software shows the same levels.

**Can I run this with GG?**  
Yeah, it works fine. Just seems pointless to have both running.

**Why did you make this?**  
The official SteelSeries software is a 500MB Electron app that wants to update constantly and show me ads. I just want to see my battery level.