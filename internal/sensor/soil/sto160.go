//go:build tinygo

package soil

import (
	"machine"

	"go.uber.org/zap"
)

// initSoilADC configures the ADC pin
func initSoilADC(log *zap.SugaredLogger, pin uint8) machine.ADC {
	adcPin := machine.Pin(pin)
	adc := machine.ADC{Pin: adcPin}
	adc.Configure(machine.ADCConfig{})
	log.Infof("soil ADC initialized: pin %d", pin)

	return adc
}

// readSoilADC reads analog moisture value from ADC pin and converts to percentage
// ESP32 ADC is 12-bit (0-4095), 3.3V reference
// Dry soil (high resistance): ~3.3V → ~4095
// Wet soil (low resistance): ~0V → ~0
// Returns percentage (0% = wet, 100% = dry)
func readSoilADC(log *zap.SugaredLogger, adc machine.ADC) (*float32, error) {
	rawValue := uint32(adc.Get())

	// convert to percentage: normalize to 0-100% range
	percentage := float32(rawValue) / 4095.0 * 100.0

	log.Infof("soil ADC reading: raw=%d, moisture=%.2f%%", rawValue, percentage)
	return &percentage, nil
}
