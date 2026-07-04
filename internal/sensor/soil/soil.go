package soil

import (
	"time"

	"go.uber.org/zap"
)

type Reading struct {
	Moisture *float32
}

// periodically read soil sensor and sends data to channel
func Start(log *zap.SugaredLogger, pin uint8, interval time.Duration, readingChan chan<- Reading) {

	// initialize the soil sensor ADC pin
	adc := initSoilADC(log, pin)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		moisture, err := readSoilADC(log, adc)
		if err != nil {
			log.Errorf("soil sensor error: %v", err)
		}
		readingChan <- Reading{
			Moisture: moisture,
		}
	}
}
