//go:build tinygo

package climate

import (
	"machine"

	"go.uber.org/zap"
	"tinygo.org/x/drivers/dht"
)

var (
	dht22Sensor dht.Device
)

// initDHT22 initializes the DHT22 sensor on the specified pin
func initDHT22(log *zap.SugaredLogger, pin int) {
	dht22Sensor = dht.New(machine.Pin(pin), dht.DHT22)
	log.Infof("DHT22 initialized on pin %d", pin)
}

// readDHT22 reads temperature and humidity from DHT22 sensor using official TinyGo driver
// Returns (temperature, humidity, error)
// Temperature in °C, Humidity in %
func readDHT22(log *zap.SugaredLogger) (*float32, *float32, error) {
	err := dht22Sensor.ReadMeasurements()
	if err != nil {
		log.Errorf("DHT22 read error: %v", err)
		return nil, nil, err
	}

	temp, humidity, err := dht22Sensor.Measurements()
	if err != nil {
		log.Errorf("DHT22 measurements error: %v", err)
		return nil, nil, err
	}

	// Convert from int16/uint16 (1/10 units) to float32 (actual values)
	// temp: 350 = 35.0°C, humidity: 652 = 65.2%
	tempFloat := float32(temp) / 10.0
	humiFloat := float32(humidity) / 10.0

	log.Infof("climate reading: %.1f°C, %.1f%% humidity", tempFloat, humiFloat)

	return &tempFloat, &humiFloat, nil
}
