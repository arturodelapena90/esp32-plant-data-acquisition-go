package mqtt

import (
	"context"
	"encoding/json"
	"io"
	"time"

	mqtt "github.com/soypat/natiu-mqtt"
	"go.uber.org/zap"
)

// Data structure matches your existing payload
type Data struct {
	Timestamp   int64    `json:"timestamp"`
	Light       *float32 `json:"light_lux"`
	Temperature *float32 `json:"temperature_c"`
	Humidity    *float32 `json:"humidity_percent"`
	Moisture1   *float32 `json:"moisture1_percent"`
	Moisture2   *float32 `json:"moisture2_percent"`
}

// Client wraps the natiu-mqtt client
type Client struct {
	client *mqtt.Client
}

// SetupMQTT initializes the client.
func SetupMQTT(conn io.ReadWriteCloser, clientID string) (*Client, error) {
	client := mqtt.NewClient(mqtt.ClientConfig{
		Decoder: mqtt.DecoderNoAlloc{make([]byte, 1500)},
	})

	var varconn mqtt.VariablesConnect
	varconn.SetDefaultMQTT([]byte(clientID))
	varconn.KeepAlive = 60

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Connect(ctx, conn, &varconn)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// Publish serializes and sends your data
func (c *Client) Publish(log *zap.SugaredLogger, topic string, mqttChan <-chan Data) error {
	for data := range mqttChan {
		payload, err := json.Marshal(data)
		if err != nil {
			log.Errorf("failed to marshal data: %v", err)
			continue
		}
		pubFlags, _ := mqtt.NewPublishFlags(mqtt.QoS0, false, false)
		pubVar := mqtt.VariablesPublish{TopicName: []byte(topic)}
		if err := c.client.PublishPayload(pubFlags, pubVar, payload); err != nil {
			log.Errorf("failed to publish data: %v", err)
		}
	}
	return nil
}
