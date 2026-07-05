//go:build tinygo

package climate

import (
	"machine"

	"go.uber.org/zap"
	"tinygo.org/x/drivers/dht"
)

func initDHT22(log *zap.SugaredLogger, pin uint8) (dht.Device, error) {
	device := dht.New(machine.Pin(pin), dht.DHT22)
	return device, nil
}

func readDHT22(log *zap.SugaredLogger, device dht.Device) (*float32, *float32, error) {
	if err := device.ReadMeasurements(); err != nil {
		log.Errorf("DHT22 read error: %v", err)
		return nil, nil, err
	}

	temp, humidity, err := device.Measurements()
	if err != nil {
		log.Errorf("DHT22 measurements error: %v", err)
		return nil, nil, err
	}

	tempFloat := float32(temp) / 10.0
	humiFloat := float32(humidity) / 10.0

	log.Infof(
		"climate reading: %.1f°C, %.1f%% humidity",
		tempFloat,
		humiFloat,
	)

	return &tempFloat, &humiFloat, nil
}
