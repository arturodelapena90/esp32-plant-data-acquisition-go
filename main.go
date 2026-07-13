package main

import (
	"device/esp"
	"fmt"
	"machine"
	"net"
	"time"

	"tinygo.org/x/drivers/netdev"
	nl "tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/aggregator"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/config"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/mqtt"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/light"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/soil"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/statusled"
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
	// Status LED Setup
	// --------------------
	cfg.StatusLEDPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	statusled.SetColor(cfg.StatusLEDPin, 0, 0, 0)
	fail := func(err error) {
		fmt.Println("fatal:", err)
		for i := 0; i < 3; i++ {
			statusled.SetColor(cfg.StatusLEDPin, 20, 0, 0)
		}
		time.Sleep(20 * time.Second)
		esp.RTC_CNTL.SetOPTIONS0_SW_SYS_RST(1)
		for {
		} // unreachable once the reset register takes effect
	}

	// --------------------
	// WiFi Setup
	// --------------------
	wifi := link.Esplink{}
	netdev.UseNetdev(&wifi)

	if err := wifi.NetConnect(&nl.ConnectParams{
		Ssid:       cfg.WifiSSID,
		Passphrase: cfg.WifiPassword,
	}); err != nil {
		fail(fmt.Errorf("failed to connect to WiFi: %w", err))
	}
	fmt.Printf("WiFi connected: %s\n", cfg.WifiSSID)

	// --------------------
	// Establish TCP Pipe
	// --------------------
	conn, err := net.Dial("tcp", cfg.MQTTBroker)
	if err != nil {
		fail(fmt.Errorf("failed to connect to MQTT broker: %w", err))
	}
	defer conn.Close()

	// --------------------
	// MQTT Setup
	// --------------------
	mqttClient, err := mqtt.SetupMQTT(conn, cfg.MQTTClientID)
	if err != nil {
		fail(fmt.Errorf("MQTT init failed: %w", err))
	}

	fmt.Printf("MQTT broker connected: %s\n", cfg.MQTTBroker)

	// --------------------
	// I2C setup
	// --------------------
	bus := machine.I2C1
	bus.Configure(machine.I2CConfig{
		SDA: cfg.I2CSDAPin,
		SCL: cfg.I2CSCLPin,
	})

	// --------------------
	// Sensors
	// --------------------
	lightSensor, err := light.New(bus, 0x23)
	if err != nil {
		fail(fmt.Errorf("light init failed: %w", err))
	}

	climateSensor, err := climate.New(cfg.DHT22Pin)
	if err != nil {
		fail(fmt.Errorf("climate init failed: %w", err))
	}

	soil1, err := soil.New(cfg.SoilPin1)
	if err != nil {
		fail(fmt.Errorf("soil1 init failed: %w", err))
	}

	soil2, err := soil.New(cfg.SoilPin2)
	if err != nil {
		fail(fmt.Errorf("soil2 init failed: %w", err))
	}

	// --------------------
	// Channels
	// --------------------
	lightChan := make(chan light.Reading)
	climateChan := make(chan climate.Reading)
	soilChan1 := make(chan soil.Reading)
	soilChan2 := make(chan soil.Reading)
	mqttChan := make(chan mqtt.Data)
	aggChan := make(chan mqtt.Data)

	// --------------------
	// Start sensors
	// --------------------
	go lightSensor.Start(cfg.ReadInterval, lightChan)
	go climateSensor.Start(cfg.ReadInterval, climateChan)
	go soil1.Start(cfg.ReadInterval, soilChan1)
	go soil2.Start(cfg.ReadInterval, soilChan2)

	statusled.SetColor(cfg.StatusLEDPin, 0, 20, 0)

	// --------------------
	// Pipeline
	// --------------------
	go aggregator.Start(lightChan, climateChan, soilChan1, soilChan2, aggChan)
	go mqttClient.Publish(cfg.MQTTTopic, mqttChan)

	go func() {
		for data := range aggChan {
			statusled.SetColor(cfg.StatusLEDPin, 15, 0, 20)
			time.Sleep(time.Second)
			statusled.SetColor(cfg.StatusLEDPin, 0, 20, 0)

			select {
			case mqttChan <- data:
			default:
			}
		}
	}()

	fmt.Println("ESP32 Plant Data Acquisition started")

	select {}
}
