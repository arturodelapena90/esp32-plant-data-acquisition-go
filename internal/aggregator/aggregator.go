package aggregator

import (
	"time"

	"go.uber.org/zap"

	"github.com/arturodelapena90/esp32-plant-acquisition/internal/mqtt"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/climate"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/light"
	"github.com/arturodelapena90/esp32-plant-acquisition/internal/sensor/soil"
)

// begin the aggregation goroutine
func Start(log *zap.SugaredLogger, lightChan <-chan light.Reading, climateChan <-chan climate.Reading, soilChan1 <-chan soil.Reading, soilChan2 <-chan soil.Reading, mqttChan chan<- mqtt.Data) {
	for {
		// wait for readings from all sensors
		lightReading := <-lightChan
		climateReading := <-climateChan
		soilReading1 := <-soilChan1
		soilReading2 := <-soilChan2

		payload := mqtt.Data{
			Timestamp:   time.Now().Unix(),
			Light:       lightReading.Lux,
			Temperature: climateReading.Temperature,
			Humidity:    climateReading.Humidity,
			Moisture1:   soilReading1.Moisture,
			Moisture2:   soilReading2.Moisture,
		}

		log.Infof("aggregator: ts=%d light=%v temp=%v humidity=%v moisture1=%v moisture2=%v", payload.Timestamp, payload.Light, payload.Temperature, payload.Humidity, payload.Moisture1, payload.Moisture2)
		mqttChan <- payload
	}
}
