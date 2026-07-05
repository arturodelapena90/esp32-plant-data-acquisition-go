package light

import (
	"time"

	"machine"

	"go.uber.org/zap"
)

type Sensor struct {
	log  *zap.SugaredLogger
	bus  machine.I2C
	addr uint8
}

type Reading struct {
	Lux *float32
}

// New initializes the BH1750 sensor
func New(log *zap.SugaredLogger, bus machine.I2C, addr uint8) (*Sensor, error) {
	if err := initBH1750(bus, addr); err != nil {
		return nil, err
	}

	return &Sensor{
		log:  log,
		bus:  bus,
		addr: addr,
	}, nil
}

// Start continuously reads light data and sends it to channel
func (s *Sensor) Start(interval time.Duration, out chan<- Reading) {
	// immediate first reading
	if r, err := s.Read(); err == nil {
		out <- r
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		r, err := s.Read()
		if err != nil {
			s.log.Errorf("light sensor error: %v", err)
			continue
		}
		out <- r
	}
}

// Read performs a single lux measurement
func (s *Sensor) Read() (Reading, error) {
	lux, err := readBH1750(s.bus, s.addr)
	return Reading{Lux: lux}, err
}
