package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {

	// WiFi
	WifiSSID     string `env:"WIFI_SSID,required"`
	WifiPassword string `env:"WIFI_PASSWORD,required"`

	// MQTT
	MQTTBroker string `env:"MQTT_BROKER,required"`
	MQTTTopic  string `env:"MQTT_TOPIC,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
