//go:build tinygo

package light

import (
	"machine"
)

const (
	BH1750PowerOn  = 0x01
	BH1750Reset    = 0x07
	BH1750ModeHRes = 0x10
)

// initBH1750 configures the sensor
func initBH1750(bus machine.I2C, addr uint8) error {
	// power on
	if err := bus.Tx(addr, []byte{BH1750PowerOn}, nil); err != nil {
		return err
	}

	// reset
	if err := bus.Tx(addr, []byte{BH1750Reset}, nil); err != nil {
		return err
	}

	// set mode
	if err := bus.Tx(addr, []byte{BH1750ModeHRes}, nil); err != nil {
		return err
	}

	return nil
}

// readBH1750 reads lux from sensor
func readBH1750(bus machine.I2C, addr uint8) (*float32, error) {
	data := make([]byte, 2)

	if err := bus.Tx(addr, nil, data); err != nil {
		return nil, err
	}

	raw := uint16(data[0])<<8 | uint16(data[1])
	lux := float32(raw) / 1.2

	return &lux, nil
}
