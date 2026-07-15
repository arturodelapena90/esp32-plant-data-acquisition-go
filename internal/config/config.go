package config

import (
	"fmt"
	"machine"
	"strconv"
	"time"
)

// envs are injected at build time via `-ldflags -X` in the makefile
var (
	buildWifiSSID         string
	buildWifiPassword     string
	buildRaspberryPiIP    string
	buildMQTTTopic        string
	buildMQTTClientID     string
	buildMQTTPort         string
	buildReadIntervalSecs string
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
	if buildWifiSSID == "" || buildWifiPassword == "" || buildRaspberryPiIP == "" || buildMQTTTopic == "" ||
		buildMQTTClientID == "" || buildMQTTPort == "" || buildReadIntervalSecs == "" {
		return nil, fmt.Errorf("missing required build-time config (wifi/broker credentials) — build with `make build`/`make flash`, not a bare tinygo build")
	}

	readIntervalSecs, err := strconv.Atoi(buildReadIntervalSecs)
	if err != nil || readIntervalSecs <= 0 {
		return nil, fmt.Errorf("invalid READ_INTERVAL in .env: must be a positive integer number of seconds, got %q", buildReadIntervalSecs)
	}

	cfg := &Config{
		RaspberryPiIP: buildRaspberryPiIP,
		WifiSSID:      buildWifiSSID,
		WifiPassword:  buildWifiPassword,
		MQTTTopic:     buildMQTTTopic,
		MQTTClientID:  buildMQTTClientID,
		DHT22Pin:      4,
		SoilPin1:      7,
		SoilPin2:      6,
		I2CSDAPin:     8,
		I2CSCLPin:     9,
		ReadInterval:  time.Duration(readIntervalSecs) * time.Second,
	}
	cfg.MQTTBroker = cfg.RaspberryPiIP + ":" + buildMQTTPort
	return cfg, nil
}
