package climate

import (
	"time"

	"go.uber.org/zap"
)

type Reading struct {
	Temperature *float32
	Humidity    *float32
}

// Init initializes the DHT22 sensor hardware on the specified pin
func Init(log *zap.SugaredLogger, pin int) {
	initDHT22(log, pin)
}

// periodically read temperature and humidity from DHT22 sensor and sends data to channel
func Read(log *zap.SugaredLogger, interval time.Duration, readingChan chan<- Reading) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		temp, humidity, err := readDHT22(log)
		if err != nil {
			log.Errorf("climate sensor error: %v", err)
		}
		readingChan <- Reading{
			Temperature: temp,
			Humidity:    humidity,
		}
	}
}
