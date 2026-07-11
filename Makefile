TARGET := esp32s3-generic
CONFIG_PKG := github.com/arturodelapena90/esp32-plant-acquisition/internal/config
SERIAL_PORT ?= /dev/ttyUSB0

# .env has no effect on the flashed device (bare-metal TinyGo has no OS
# environment for os.Getenv to read at runtime) — it's only a source of
# truth for these -ldflags, built here at compile time. See config.go.
LDFLAGS = $(shell set -a && . ./.env && set +a && \
	echo "-X $(CONFIG_PKG).buildWifiSSID=$$WIFI_SSID" \
	     "-X $(CONFIG_PKG).buildWifiPassword=$$WIFI_PASSWORD" \
	     "-X $(CONFIG_PKG).buildRaspberryPiIP=$$RASPBERRY_PI_IP" \
	     "-X $(CONFIG_PKG).buildMQTTTopic=$$MQTT_TOPIC" \
	     "-X $(CONFIG_PKG).buildMQTTClientID=$$MQTT_CLIENT_ID")

.PHONY: build flash monitor

build:
	tinygo build -target=$(TARGET) -ldflags="$(LDFLAGS)" -o firmware.bin main.go

flash:
	tinygo flash -target=$(TARGET) -ldflags="$(LDFLAGS)" main.go

monitor:
	tinygo monitor -port=$(SERIAL_PORT) -baud=921600
