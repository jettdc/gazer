package config

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
	ColorCode        string `yaml:"color,omitempty" structs:"-"`
	Retries          uint8  `yaml:"retries,omitempty" structs:"-"`
	AttemptedRetries uint8  `yaml:"-" structs:"-"`
}

const (
	RestartAlways uint8 = iota
	RestartNo
	RestartUnlessStopped
	RestartOnFailure
)

const (
	Reset string = "\033[0m"

	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)
