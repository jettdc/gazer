package main

import (
	"context"
	"fmt"
	"github.com/jettdc/gazer/config"
	"github.com/jettdc/gazer/gaze"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	c, err := config.LoadConfig("./gazer.yaml", ctx)
	if err != nil {
		panic(err)
	}

	go gaze.ExecuteDescriptions(c.Shell, c.Descriptions)

	restartRequests := make(chan string, 1)
	go gaze.Watch(c.Descriptions, restartRequests)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	go func() {
		for {
			select {
			case restartRequest := <-restartRequests:

				for i, _ := range *c.Descriptions {
					desc := &(*c.Descriptions)[i]
					if desc.Name == restartRequest {
						fmt.Println(desc.Output("Restarting process..."))

						desc.Cancel()
						desc.Context, desc.Cancel = context.WithCancel(ctx)

						descs := &[]config.GazeDesc{
							*desc,
						}

						go gaze.ExecuteDescriptions(c.Shell, descs)
						break
					}
				}

			case <-sigs:
				for i, _ := range *c.Descriptions {
					desc := &(*c.Descriptions)[i]
					desc.Cancel()
				}
				os.Exit(0)
			}
		}
	}()

	for {
		var w1, w2 string
		_, err := fmt.Scanln(&w1, &w2)
		if err != nil {
			fmt.Println("[gazer] Invalid command. Try \"rs <process name>\"")
			continue
		}

		if w1 != "rs" {
			fmt.Println("[gazer] Invalid command. Try \"rs <process name>\"")
			continue
		}

		found := false
		for i, _ := range *c.Descriptions {
			desc := &(*c.Descriptions)[i]
			if desc.Name == w2 {
				restartRequests <- w2
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("[gazer] Process with name %s not found.\n", w2)
		}
	}
}
