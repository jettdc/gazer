package gaze

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jettdc/gazer/config"
	"sync"
	"time"
)

func Watch(c *[]config.GazeDesc, restartRequests chan string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(*c))

	for i, _ := range *c {
		desc := &(*c)[i]
		go watchForDescription(desc, restartRequests, wg)
	}

	wg.Wait()
}

func watchForDescription(desc *config.GazeDesc, restartRequests chan string, wg *sync.WaitGroup) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		lastEvent := time.Now()
		for {
			select {
			case <-watcher.Events:
				// Debounce
				if time.Since(lastEvent) >= 500*time.Millisecond {
					restartRequests <- desc.Name
					lastEvent = time.Now()
				}

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	for _, toWatch := range desc.WatchPaths {
		if err := watcher.Add(toWatch); err != nil {
			fmt.Println("ERROR", err)
		}
	}

	<-done
	wg.Done()
}
