package main

import (
	"machine"
	"net"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/aggregator"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/config"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/logger"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/mqtt"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/light"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/soil"
)

func main() {

	// --------------------
	// Load ENVs
	// --------------------
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// --------------------
	// Logger Setup
	// --------------------
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Log
	log.Info("ESP32 Plant Data Acquisition started")

	// --------------------
	// WiFi Setup
	// --------------------
	if err := machine.WIFI.Configure(machine.WIFIConfig{
		SSID:     cfg.WifiSSID,
		Password: cfg.WifiPassword,
	}); err != nil {
		log.Fatalf("failed to configure WiFi: %v", err)
	}

	// --------------------
	// Establish TCP Pipe
	// --------------------
	conn, err := net.Dial("tcp", cfg.RaspberryPiIP)
	if err != nil {
		log.Fatalf("failed to connect to MQTT broker: %v", err)
	}
	defer conn.Close()

	// --------------------
	// MQTT Setup
	// --------------------
	mqttClient, err := mqtt.SetupMQTT(conn, cfg.MQTTBroker)
	if err != nil {
		log.Fatalf("MQTT Init failed: %v", err)
	}

	log.Infof("MQTT broker connected: %s", cfg.MQTTBroker)

	// --------------------
	// I2C setup
	// --------------------
	bus := machine.I2C1
	bus.Configure(machine.I2CConfig{
		SDA: machine.GP8,
		SCL: machine.GP9,
	})

	// --------------------
	// Sensors
	// --------------------
	lightSensor, err := light.New(log, bus, 0x23)
	if err != nil {
		log.Fatalf("light init failed: %v", err)
	}

	climateSensor, err := climate.New(log, cfg.DHT22Pin)
	if err != nil {
		log.Fatalf("climate init failed: %v", err)
	}

	soil1, err := soil.New(log, cfg.SoilPin1)
	if err != nil {
		log.Fatalf("soil1 init failed: %v", err)
	}

	soil2, err := soil.New(log, cfg.SoilPin2)
	if err != nil {
		log.Fatalf("soil2 init failed: %v", err)
	}

	// --------------------
	// Channels
	// --------------------
	lightChan := make(chan light.Reading)
	climateChan := make(chan climate.Reading)
	soilChan1 := make(chan soil.Reading)
	soilChan2 := make(chan soil.Reading)
	mqttChan := make(chan mqtt.Data)

	// --------------------
	// Start sensors
	// --------------------
	go lightSensor.Start(cfg.ReadInterval, lightChan)
	go climateSensor.Start(cfg.ReadInterval, climateChan)
	go soil1.Start(cfg.ReadInterval, soilChan1)
	go soil2.Start(cfg.ReadInterval, soilChan2)

	// --------------------
	// Pipeline
	// --------------------
	go aggregator.Start(log, lightChan, climateChan, soilChan1, soilChan2, mqttChan)
	go mqttClient.Publish(log, cfg.MQTTTopic, mqttChan)

	log.Info("system running")

	select {}
}
