package cpu_temp

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

	tempRegex := regexp.MustCompile(`([+-]?\d+(\.\d+)?)°?C`)

	for chip, values := range sensors.Chips {
		if strings.Contains(chip, "coretemp") || strings.Contains(chip, "amd") || strings.Contains(chip, "acpitz") {
			for key, value := range values {
				if strings.Contains(key, "Core") || strings.Contains(key, "Package") {
					matches := tempRegex.FindStringSubmatch(value)
					if len(matches) > 0 {
						return &Stats{Temperature: matches[0]}, nil
					}
				}
			}
		}
	}

	// Если температура не найдена
	return nil, errors.New("not found chip or temperature data")
}
