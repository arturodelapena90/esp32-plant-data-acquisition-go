package soil

import (
	"machine"
	"time"

	"go.uber.org/zap"
)

type Sensor struct {
	log *zap.SugaredLogger
	adc machine.ADC
}

type Reading struct {
	Moisture *float32
}

func New(log *zap.SugaredLogger, pin uint8) (*Sensor, error) {
	adc, err := initSoilADC(log, pin)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		log: log,
		adc: adc,
	}, nil
}

func (s *Sensor) Start(interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		reading, err := s.Read()
		if err != nil {
			s.log.Errorf("soil sensor error: %v", err)
		}

		readingChan <- reading
	}
}

func (s *Sensor) Read() (Reading, error) {
	moisture, err := readSoilADC(s.log, s.adc)

	return Reading{
		Moisture: moisture,
	}, err
}
