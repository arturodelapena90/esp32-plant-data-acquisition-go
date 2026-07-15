package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	mqtt "github.com/soypat/natiu-mqtt"
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

// String implements fmt.Stringer so %v/%s on a Data print readable values
// instead of pointer addresses for the nil-able *float32 fields.
func (d Data) String() string {
	return fmt.Sprintf(
		"ts=%d light=%s temp=%s humidity=%s moisture1=%s moisture2=%s",
		d.Timestamp, formatReading(d.Light), formatReading(d.Temperature),
		formatReading(d.Humidity), formatReading(d.Moisture1), formatReading(d.Moisture2),
	)
}

func formatReading(v *float32) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *v)
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
func (c *Client) Publish(topic string, mqttChan <-chan Data) error {
	for data := range mqttChan {
		payload, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("failed to marshal data: %v\n", err)
			continue
		}
		pubFlags, _ := mqtt.NewPublishFlags(mqtt.QoS0, false, false)
		pubVar := mqtt.VariablesPublish{TopicName: []byte(topic), PacketIdentifier: 1}
		if err := c.client.PublishPayload(pubFlags, pubVar, payload); err != nil {
			fmt.Printf("failed to publish data: %v\n", err)
		}
	}
	return nil
}
