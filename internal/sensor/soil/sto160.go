//go:build tinygo

package soil

import (
	"machine"

	"go.uber.org/zap"
)

func initSoilADC(log *zap.SugaredLogger, pin uint8) (machine.ADC, error) {
	adcPin := machine.Pin(pin)

	adc := machine.ADC{Pin: adcPin}
	adc.Configure(machine.ADCConfig{})

	return adc, nil
}

func readSoilADC(log *zap.SugaredLogger, adc machine.ADC) (*float32, error) {
	raw := uint32(adc.Get())

	percentage := float32(raw) / 4095 * 100

	log.Infof(
		"soil ADC reading: raw=%d moisture=%.2f%%",
		raw,
		percentage,
	)

	return &percentage, nil
}
