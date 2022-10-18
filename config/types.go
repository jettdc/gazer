package config

import (
	"context"
	"fmt"
	"github.com/jettdc/gazer/out"
)

type Config struct {
	Shell        *ShellConfig `yaml:"shell"`
	Descriptions *[]GazeDesc  `yaml:"gaze-at"`
}

type ShellConfig struct {
	Executable string   `yaml:"exec"`
	Arguments  []string `yaml:"args,omitempty" structs:"-"`
}

type GazeDesc struct {
	Name    string `yaml:"name"`
	Command string `yaml:"cmd"`
	//Watch         bool
	//RestartPolicy uint8
	ColorCode        string             `yaml:"color,omitempty" structs:"-"`
	Restart          string             `yaml:"restart,omitempty" structs:"-"`
	Retries          int32              `yaml:"retries,omitempty" structs:"-"`
	AttemptedRetries int32              `yaml:"-" structs:"-"`
	WatchPaths       []string           `yaml:"watch,omitempty" structs:"-"`
	Context          context.Context    `yaml:"-" structs:"-"`
	Cancel           context.CancelFunc `yaml:"-" structs:"-"`
}

func (description *GazeDesc) Output(text string) string {
	return out.ColorText(fmt.Sprintf("[%s] %s", description.Name, text), description.ColorCode)
}

func (description *GazeDesc) ErrorOutput(text string) string {
	return out.ColorText(fmt.Sprintf("[%s | ERROR] %s", description.Name, text), description.ColorCode)
}

func (description *GazeDesc) AdminOutput(text string) string {
	return out.ColorText(fmt.Sprintf("[gazer] %s", text), description.ColorCode)
}
