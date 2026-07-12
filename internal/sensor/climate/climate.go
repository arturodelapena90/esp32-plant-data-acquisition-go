package climate

import (
	"machine"
	"time"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate/dhtdriver"
)

type Sensor struct {
	device dhtdriver.Device
}

type Reading struct {
	Temperature *float32
	Humidity    *float32
}

func New(pin machine.Pin) (*Sensor, error) {
	device, err := initDHT22(pin)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		device: device,
	}, nil
}

func (s *Sensor) Start(interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		reading, _ := s.Read()
		readingChan <- reading
	}
}

func (s *Sensor) Read() (Reading, error) {
	temp, humidity, err := readDHT22(s.device)

	return Reading{
		Temperature: temp,
		Humidity:    humidity,
	}, err
}
