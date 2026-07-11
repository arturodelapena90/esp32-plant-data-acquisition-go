# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Go firmware for an ESP32-S3, compiled with TinyGo, that reads sensors on a habanero plant (BH1750 light over I2C, DHT22 temperature/humidity over GPIO, 2x capacitive soil moisture over ADC) and publishes JSON readings to an MQTT broker over WiFi. The broker and WiFi AP are both hosted on a Raspberry Pi.

## Build / flash / monitor

This is TinyGo firmware, not a regular Go binary — `go build` will fail because it imports the `machine` package, which only resolves under the TinyGo toolchain. Both `go` (1.24.1) and `tinygo` (0.41.1) are installed.

```bash
# Build
tinygo build -target=esp32-s3 -o firmware.bin main.go

# Flash
tinygo flash -target=esp32-s3 main.go

# Monitor serial output
tinygo monitor -port=/dev/ttyUSB0 -baud=921600
```

Use `go vet ./...` / `gofmt` for quick static checks on non-hardware code, but treat `tinygo build` as the real compile check since hardware-specific files are gated behind `//go:build tinygo` (see below).

## Configuration

Runtime config comes from environment variables (`.env`, loaded via `github.com/caarlos0/env/v10` in `internal/config/config.go`), not from flashed-in constants:

- `RASPBERRY_PI_IP` — Pi's IP; MQTT broker URL is derived as `mqtt://<ip>:1883`
- `WIFI_SSID`, `WIFI_PASSWORD` — Pi's AP credentials
- `MQTT_TOPIC` — publish topic

Hardware pin assignments (`DHT22Pin`, `SoilPin1`, `SoilPin2`) and `ReadInterval` are hardcoded defaults in `config.LoadConfig()`, not env-driven. `.env` is gitignored — don't commit real credentials into it or into source.

## Architecture

Each sensor type lives in its own package under `internal/sensor/` and follows the same shape:

- A platform-independent file (e.g. `light.go`, `climate.go`, `soil.go`) defining `Sensor`, `Reading`, `New(...)`, `Start(interval, chan<- Reading)`, and `Read()`.
- A `//go:build tinygo`-gated file (e.g. `bh1750.go`, `dht22.go`, `sto160.go`) with the actual register/protocol-level driver code. This split exists so the non-hardware logic can at least be parsed/vetted with plain `go`, while the real hardware access only builds under TinyGo.

Data flow, wired up in `main.go`:

1. `config.LoadConfig()` reads env vars, then WiFi is configured and a TCP connection to the Pi is dialed.
2. `mqtt.SetupMQTT()` (wrapping `github.com/soypat/natiu-mqtt`) connects over that TCP conn.
3. I2C bus (`machine.I2C1` on GP8/GP9) is configured once and shared; each sensor's `New()` is called to initialize it.
4. Each sensor runs `Start()` in its own goroutine, ticking on `cfg.ReadInterval` and pushing a `Reading` onto its own channel (light immediately reads once before the first tick; climate/soil wait for the first tick).
5. `aggregator.Start()` runs in its own goroutine: it blocks on receiving one reading from each of the four channels (light, climate, soil1, soil2), merges them into an `mqtt.Data` payload with a fresh Unix timestamp, and pushes it to `mqttChan`. Because it waits on all four channels every loop, publish rate is effectively gated by the slowest sensor.
6. `mqttClient.Publish()` drains `mqttChan`, JSON-marshals each `mqtt.Data`, and publishes it (QoS0) to `cfg.MQTTTopic`.

All sensor `Reading` fields are `*float32` pointers (nil on read error), and `mqtt.Data` mirrors that — a nil field means that sensor failed on that cycle but the payload is still published with the other readings.

Logging is a single global `*zap.SugaredLogger` (`logger.Log`), initialized once in `main.go` via `logger.Init()` and threaded explicitly into every constructor/`Start`/`Read` call rather than accessed as a package global from the callees.

Note: `README.md`'s "Architecture" section describes an older single-file `main.go` layout (`readLightSensor()`, `aggregateReadings()`, etc.) that predates the current `internal/` package split — trust the code over that section.
