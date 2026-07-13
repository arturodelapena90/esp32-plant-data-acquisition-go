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

// envs are injected at build time via `-ldflags -X` in the makefile
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
	StatusLEDPin machine.Pin
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
		StatusLEDPin:  48,
		ReadInterval:  10 * time.Second,
	}
	cfg.MQTTBroker = cfg.RaspberryPiIP + ":" + mqttPort
	return cfg, nil
}
