package cpu_temp

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ssimunic/gosensors"
)

type Stats struct {
	Temperature string
}

func Get() (*Stats, error) {
	sensors, err := gosensors.NewFromSystem()
	if err != nil {
		return nil, fmt.Errorf("error getting sensors data: %v", err)
	}

	tempRegex := regexp.MustCompile(`([+-]?\d+(\.\d+)?)°C`)

	for chip, values := range sensors.Chips {
		if chip == "coretemp-isa-0000" {
			for key, value := range values {
				if key == "Core 0" || key == "Package id 0" {
					matches := tempRegex.FindStringSubmatch(value)
					if len(matches) > 0 {
						return &Stats{Temperature: matches[0]}, nil
					}
				}
			}
		}
	}

	return nil, errors.New("not found chip or temperature data")
}
