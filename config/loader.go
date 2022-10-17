package config

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/jettdc/gazer/out"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find specified config file at %s", path)
	}
	defer file.Close()

	config := &Config{}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("invalid yaml in config file: %s", err.Error())
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	adaptConfig(config)
	return config, nil
}

func adaptConfig(config *Config) {
	// Change color names to color codes
	for i, _ := range *config.Descriptions {
		desc := &(*config.Descriptions)[i]
		if desc.ColorCode == "" {
			desc.ColorCode = out.UseNextAvailableColorCode()
		} else {
			desc.ColorCode, _ = out.UseColor(desc.ColorCode) // Already validated
		}
	}
}

func validateConfig(config *Config) error {
	if err := validateShellConfig(config.Shell); err != nil {
		return err
	}

	// Validate descriptions
	for _, desc := range *config.Descriptions {
		if err := validateDesc(&desc); err != nil {
			return err
		}
	}

	return nil
}

func validateShellConfig(shellConfig *ShellConfig) error {
	if shellConfig == nil {
		return fmt.Errorf("missing required config \"shell\"")
	}
	if structs.HasZero(shellConfig) {
		return fmt.Errorf("this sucks!")
	}
	return nil
}

func validateDesc(description *GazeDesc) error {
	if structs.HasZero(description) {
		return fmt.Errorf("missing required field in description: %s", description.Name)
	}

	if description.ColorCode != "" {
		_, err := out.GetColorByName(description.ColorCode) // not currently a code, just a name
		if err != nil {
			return err
		}
	}
	return nil
}

func getCmd(key string, description map[string]interface{}) (string, error) {
	cmd, ok := description["cmd"]
	if !ok {
		return "", fmt.Errorf("invalid config: specification for %s missing required parameter \"cmd\"", key)
	}

	cmd, ok = cmd.(string)
	if !ok {
		return "", fmt.Errorf("invalid config: \"cmd\" for %s must be a string", key)
	}

	return cmd.(string), nil
}

func getColorCode(key string, description map[string]interface{}) (string, error) {
	color, ok := description["color"]
	if !ok {
		return out.UseNextAvailableColorCode(), nil
	}

	color, ok = color.(string)
	if !ok {
		return "", fmt.Errorf("invalid config: \"color\" for %s must be a string", key)
	}

	code, err := out.UseColor(color.(string))
	if err != nil {
		return "", fmt.Errorf("invalid config: \"color\" for %s is not an available option.", key)
	}

	return code, nil
}
