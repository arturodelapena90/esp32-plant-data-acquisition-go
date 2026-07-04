package main

import (
	"time"

	mqttclient "github.com/eclipse/paho.mqtt.golang"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/aggregator"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/logger"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/mqtt"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/light"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/soil"
)

// configuration
const (
	DHT22Pin     = 4
	SoilPin1     = 7
	SoilPin2     = 8
	ReadInterval = 30 * time.Second
	MQTTBroker   = "mqtt://192.168.1.100:1883"
	MQTTTopic    = "esp32/habanero-plant/data"
)

func main() {
	// initialize logger
	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Log
	log.Info("ESP32 Plant Data Acquisition started")

	// setup MQTT connection
	opts := mqttclient.NewClientOptions()
	opts.AddBroker(MQTTBroker)
	opts.SetClientID("esp32-habanero")
	opts.SetAutoReconnect(true)

	client := mqttclient.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to connect to MQTT broker: %v", token.Error())
	}
	log.Infof("MQTT broker connected: %s", MQTTBroker)
	defer client.Disconnect(250)

	// TODO: remove inits for the rest of the sensors, handle in the sensor package
	// initialize sensors
	light.Init(log)
	climate.Init(log, DHT22Pin)

	// create data channels
	lightChan := make(chan light.Reading)
	climateChan := make(chan climate.Reading)
	soilChan1 := make(chan soil.Reading)
	soilChan2 := make(chan soil.Reading)
	mqttChan := make(chan mqtt.Data)

	// start sensor goroutines
	go light.Read(log, ReadInterval, lightChan)
	go climate.Read(log, ReadInterval, climateChan)
	go soil.Start(log, SoilPin1, ReadInterval, soilChan1)
	go soil.Start(log, SoilPin2, ReadInterval, soilChan2)

	// start aggregator
	go aggregator.Start(log, lightChan, climateChan, soilChan1, soilChan2, mqttChan)

	// start MQTT publisher
	go mqtt.Publish(log, mqttChan, client, MQTTTopic)

	log.Info("system running")

	// keep main alive
	select {}
}
