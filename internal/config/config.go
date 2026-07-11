package config

import (
	"fmt"
	"machine"
	"time"
)

const (
	mqttPort            = "1883"
	defaultMQTTClientID = "esp32-habanero-01"
)

// buildWifiSSID, buildWifiPassword, buildRaspberryPiIP, buildMQTTTopic and
// buildMQTTClientID are injected at build time via `-ldflags -X`, e.g.:
//
//	tinygo flash -target=esp32s3-generic -ldflags="-X github.com/arturodelapena90/esp32-plant-acquisition/internal/config.buildWifiSSID=..."
//
// `make build`/`make flash` do this for you by reading .env (see Makefile).
// A flashed ESP32-S3 has no OS environment for os.Getenv to read from, so
// .env can't be loaded at runtime the way it could on a hosted Go program.
//
// These must stay zero-value here: TinyGo's `-X` silently fails to override
// a var that already has a non-empty initializer (verified against TinyGo
// 0.41.1 — the linker reports no error, it just keeps the compiled-in
// default). Any default value goes in LoadConfig() instead.
var (
	buildWifiSSID      string
	buildWifiPassword  string
	buildRaspberryPiIP string
	buildMQTTTopic     string
	buildMQTTClientID  string
)

type Config struct {

	// Raspberry Pi
	RaspberryPiIP string

	// WiFi
	WifiSSID     string
	WifiPassword string

	// MQTT
	MQTTTopic    string
	MQTTClientID string
	MQTTBroker   string

	// Hardware Settings
	DHT22Pin     machine.Pin
	SoilPin1     machine.Pin
	SoilPin2     machine.Pin
	I2CSDAPin    machine.Pin
	I2CSCLPin    machine.Pin
	ReadInterval time.Duration
}

func LoadConfig() (*Config, error) {
	if buildWifiSSID == "" || buildWifiPassword == "" || buildRaspberryPiIP == "" || buildMQTTTopic == "" {
		return nil, fmt.Errorf("missing required build-time config (wifi/broker credentials) — build with `make build`/`make flash`, not a bare tinygo build")
	}

	clientID := buildMQTTClientID
	if clientID == "" {
		clientID = defaultMQTTClientID
	}

	cfg := &Config{
		RaspberryPiIP: buildRaspberryPiIP,
		WifiSSID:      buildWifiSSID,
		WifiPassword:  buildWifiPassword,
		MQTTTopic:     buildMQTTTopic,
		MQTTClientID:  clientID,
		DHT22Pin:      4,
		SoilPin1:      7,
		SoilPin2:      6,
		I2CSDAPin:     8,
		I2CSCLPin:     9,
		ReadInterval:  30 * time.Second,
	}
	cfg.MQTTBroker = cfg.RaspberryPiIP + ":" + mqttPort
	return cfg, nil
}
