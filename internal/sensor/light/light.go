package light

import (
	"time"

	"machine"
)

type Sensor struct {
	bus  *machine.I2C
	addr uint16
}

type Reading struct {
	Lux *float32
}

// New initializes the BH1750 sensor
func New(bus *machine.I2C, addr uint16) (*Sensor, error) {
	if err := initBH1750(bus, addr); err != nil {
		return nil, err
	}

	return &Sensor{
		bus:  bus,
		addr: addr,
	}, nil
}

// Start continuously reads light data and sends it to channel
func (s *Sensor) Start(interval time.Duration, out chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		r, _ := s.Read()
		out <- r
	}
}

// Read performs a single lux measurement
func (s *Sensor) Read() (Reading, error) {
	lux, err := readBH1750(s.bus, s.addr)
	return Reading{Lux: lux}, err
}
