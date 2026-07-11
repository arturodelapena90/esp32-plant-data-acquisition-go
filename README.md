# ESP32 S3 Habanero Plant Data Acquisition

Simple, clean Go system to acquire sensor data from an Habanero plant and publish via MQTT.

## Sensors

- **BH1750** - Light intensity (I2C)
- **DHT22** - Temperature & Humidity (GPIO)
- **2x Capacitive Soil Moisture** - Redundant moisture sensors (ADC)

## Architecture

```
main.go:
├── connect WiFi (tinygo.org/x/espradio)
├── dial + connect MQTT (github.com/soypat/natiu-mqtt)
├── configure shared I2C bus
├── light.Start()    → sends light.Reading to a channel
├── climate.Start()  → sends climate.Reading to a channel
├── soil.Start() x2  → sends soil.Reading to a channel each
├── aggregator.Start() → merges all readings into mqtt.Data
└── mqttClient.Publish() → publishes mqtt.Data to the broker
```

Each sensor runs in its own goroutine (`internal/sensor/{light,climate,soil}`). The aggregator waits for a reading from each sensor, combines them into `mqtt.Data`, and sends it to the MQTT publisher. See `CLAUDE.md` for the full data-flow writeup.

## Configuration

WiFi credentials, Pi IP, and MQTT topic/client ID live in `.env`, but **`.env` is not read on the device** — a flashed ESP32-S3 has no OS environment for `os.Getenv` to read from. Instead these values are compiled directly into the binary via `-ldflags -X` (see `internal/config/config.go`), and `make build`/`make flash` read `.env` and construct those flags for you. A bare `tinygo build`/`tinygo flash` (bypassing `make`) will still compile, but the firmware will panic on boot with a clear "missing required build-time config" error rather than connecting with empty credentials.

Hardware pin assignments and read interval are defaults in `internal/config/config.go`:

```go
DHT22Pin:     4, // GPIO4
SoilPin1:     7, // GPIO7 / ADC1_CH6
SoilPin2:     6, // GPIO6 / ADC1_CH5
I2CSDAPin:    8, // GPIO8
I2CSCLPin:    9, // GPIO9
ReadInterval: 30 * time.Second,
```

## Building for TinyGo

Requires TinyGo 0.41+ (for native ESP32 WiFi support via `tinygo.org/x/espradio`).

```bash
# Build (reads .env, writes firmware.bin)
make build

# Flash (reads .env)
make flash

# Monitor
make monitor                       # defaults to /dev/ttyACM0
make monitor SERIAL_PORT=/dev/ttyUSB0
```

Uses `targets/esp32s3-uart.json` (`esp32s3-generic` + `flash-method: esp32flash`), not the default `esp32s3-generic` target directly — that default resets via the native USB-Serial/JTAG peripheral, which this board (external CH34x UART bridge) doesn't have; flashing hung on "failed to sync with ESP bootloader" until switched to the classic UART reset method.

**On WSL2**: the board isn't visible to Linux until attached via `usbipd` (Windows PowerShell, as Administrator): `usbipd list` to find its `BUSID`, `usbipd bind --busid <ID>` once, then `usbipd attach --wsl --busid <ID>` every time it's replugged. Then `ls /dev/tty*` in WSL to confirm it shows up (`/dev/ttyACM0` or `/dev/ttyUSB0` depending on the board's USB chip).

## MQTT Message

Published to `MQTT_TOPIC` (see `.env`):

```json
{
  "timestamp": 1688123456,
  "light_lux": 45000.50,
  "temperature_c": 28.5,
  "humidity_percent": 65.2,
  "moisture1_percent": 72.8,
  "moisture2_percent": 71.9
}
```
