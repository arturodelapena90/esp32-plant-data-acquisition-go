package mqtt

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// data represents the complete plant sensor payload
type Data struct {
	Timestamp   int64    `json:"timestamp"`
	Light       *float32 `json:"light_lux"`
	Temperature *float32 `json:"temperature_c"`
	Humidity    *float32 `json:"humidity_percent"`
	Moisture1   *float32 `json:"moisture1_percent"`
	Moisture2   *float32 `json:"moisture2_percent"`
}

// publish sends data to MQTT broker
func Publish(log *zap.SugaredLogger, dataChan <-chan Data, client mqtt.Client, topic string) {
	for payload := range dataChan {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Errorf("failed to marshal mqtt payload: %v", err)
			continue
		}

		token := client.Publish(topic, 1, false, data)
		token.Wait()

		if token.Error() != nil {
			log.Errorf("failed to publish to %s: %v", topic, token.Error())
			continue
		}

		log.Infof("mqtt published to %s: %s", topic, string(data))
	}
}
