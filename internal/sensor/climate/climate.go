package climate

import (
	"time"

	"go.uber.org/zap"
	"tinygo.org/x/drivers/dht"
)

type Sensor struct {
	log    *zap.SugaredLogger
	device dht.Device
}

type Reading struct {
	Temperature *float32
	Humidity    *float32
}

func New(log *zap.SugaredLogger, pin uint8) (*Sensor, error) {
	device, err := initDHT22(log, pin)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		log:    log,
		device: device,
	}, nil
}

func (s *Sensor) Start(interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		reading, err := s.Read()
		if err != nil {
			s.log.Errorf("climate sensor error: %v", err)
		}

		readingChan <- reading
	}
}

func (s *Sensor) Read() (Reading, error) {
	temp, humidity, err := readDHT22(s.log, s.device)

	return Reading{
		Temperature: temp,
		Humidity:    humidity,
	}, err
}
