//go:build tinygo

package climate

import (
	"machine"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate/dhtdriver"
)

func initDHT22(pin machine.Pin) (dhtdriver.Device, error) {
	device := dhtdriver.New(pin, dhtdriver.DHT22)
	return device, nil
}

func readDHT22(device dhtdriver.Device) (*float32, *float32, error) {
	if err := device.ReadMeasurements(); err != nil {
		return nil, nil, err
	}

	temp, humidity, err := device.Measurements()
	if err != nil {
		return nil, nil, err
	}

	tempFloat := float32(temp) / 10.0
	humiFloat := float32(humidity) / 10.0

	return &tempFloat, &humiFloat, nil
}
