package main

import (
	"github.com/jettdc/gazer/config"
	"github.com/jettdc/gazer/gaze"
)

func main() {
	c, err := config.LoadConfig("./gazer.yaml")
	if err != nil {
		panic(err)
	}

	gaze.ExecuteDescriptions(*c)
}
