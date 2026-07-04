package light

import (
	"time"

	"go.uber.org/zap"
)

type Reading struct {
	Lux *float32
}

// Init initializes the light sensor hardware
func Init(log *zap.SugaredLogger) {
	initBH1750(log)
}

// periodically read the BH1750 sensor and sends data to the channel
func Start(log *zap.SugaredLogger, interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// read light sensor
		lux, err := readBH1750(log)
		if err != nil {
			log.Errorf("light sensor error: %v", err)
		}

		// send reading to channel
		readingChan <- Reading{
			Lux: lux,
		}
	}
}
