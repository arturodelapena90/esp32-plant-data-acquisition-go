# ESP32 S3 Habanero Plant Data Acquisition

Simple, clean Go system to acquire sensor data from an Habanero plant and publish via MQTT.

## Sensors

- **BH1750** - Light intensity (I2C)
- **DHT22** - Temperature & Humidity (GPIO)
- **2x Capacitive Soil Moisture** - Redundant moisture sensors (ADC)

## Architecture

```
main.go:
├── readLightSensor()      → sends LightReading to channel
├── readClimateSensor()    → sends ClimateReading to channel
├── readSoilSensor()       → sends SoilReading to channel
├── aggregateReadings()    → reads all 3 channels, combines payload
└── publishToMQTT()        → sends PlantData to MQTT broker
```

Each sensor runs in its own goroutine. The aggregator waits for a reading from each sensor, combines them into `PlantData`, and sends to MQTT publisher.

## Configuration

Edit constants in `main.go`:

```go
const (
	DHT22Pin     = 4              // GPIO pin
	SoilPin1     = 7              // ADC pin (sensor 1)
	SoilPin2     = 8              // ADC pin (sensor 2)
	ReadInterval = 30 * time.Second
	MQTTBroker   = "mqtt://192.168.1.100:1883"
	MQTTTopic    = "esp32/habanero-plant/data"
)
```

## Building for TinyGo

```bash
# Build
tinygo build -target=esp32-s3 -o firmware.bin main.go

# Flash
tinygo flash -target=esp32-s3 main.go

# Monitor
tinygo monitor -port=/dev/ttyUSB0 -baud=921600
```

## MQTT Message

Published to `esp32/habanero-plant/data`:

```json
{
  "timestamp": 1688123456,
  "light_lux": 45000.50,
  "temperature_c": 28.5,
  "humidity_percent": 65.2,
  "moisture_percent": 72.8,
  "has_errors": false
}
```