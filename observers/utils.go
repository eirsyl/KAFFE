package observers

import (
	"time"

	"github.com/kidoman/embd/convertors/mcp3008"
)

// ReadACS712AC helps reading AC values from the ACS712 chip
func ReadACS712AC(adc *mcp3008.MCP3008, channel int) (float64, error) {
	start := time.Now()
	var (
		maxValue = 0
		minValue = 1024
	)

	for time.Now().Sub(start) < 1*time.Second {
		readValue, err := adc.AnalogValueAt(channel)
		if err != nil {
			return 0, err
		}

		// see if you have a new maxValue
		if readValue > maxValue {
			maxValue = readValue
		}
		if readValue < minValue {
			minValue = readValue
		}

	}

	voltage := float64((maxValue-minValue)*5.0) / 1024.0

	var (
		mVperAmp = 185.0
	)

	vrms := (voltage / 2.0) * 0.707
	ampsRMS := (vrms * 1000.0) / mVperAmp
	return ampsRMS, nil
}
