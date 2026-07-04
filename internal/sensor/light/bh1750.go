//go:build tinygo

package light

import (
	"machine"
	"time"

	"go.uber.org/zap"
)

const (
	BH1750Addr      = 0x23                   // 7-bit I2C address
	BH1750PowerOn   = 0x01                   // Power on command
	BH1750Reset     = 0x07                   // Reset command
	BH1750ModeHRes  = 0x10                   // Continuous high resolution mode (1 lux resolution)
	BH1750ReadDelay = 180 * time.Millisecond // Measurement time for high resolution mode
)

// initialize the BH1750 sensor via I2C
func initBH1750(log *zap.SugaredLogger) (machine.I2C, error) {
	err := machine.I2C1.Tx(BH1750Addr, []byte{BH1750PowerOn, BH1750Reset, BH1750ModeHRes}, nil)
	if err != nil {
		log.Errorf("failed to initialize BH1750: %v", err)
		return machine.I2C1, err
	}
	log.Infof("BH1750 initialized successfully")
	return machine.I2C1, nil
}

// readBH1750 reads light intensity from BH1750 sensor via I2C
func readBH1750(log *zap.SugaredLogger, i2cBus machine.I2C) (*float32, error) {

	// read 2 bytes (MSB, LSB)
	data := make([]byte, 2)
	err := i2cBus.Tx(BH1750Addr, nil, data)
	if err != nil {
		log.Errorf("failed to read BH1750 data: %v", err)
		return nil, err
	}

	// convert bytes to lux: (MSB << 8 | LSB) / 1.2
	rawValue := uint16(data[0])<<8 | uint16(data[1])
	lux := float32(rawValue) / 1.2

	return &lux, nil
}
