package config

import (
	"fmt"
	"strings"
)

func ValidateRestartParameter(parameter string) error {
	switch strings.ToLower(parameter) {
	case "retry", "always":
		return nil
	default:
		return fmt.Errorf("invalid retry parameter: %s", parameter)
	}
}
