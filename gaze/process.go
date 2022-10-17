package gaze

import (
	"bufio"
	"fmt"
	"github.com/jettdc/gazer/config"
	"github.com/jettdc/gazer/out"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func ExecuteDescriptions(c config.Config) error {
	output := make(chan string, 8)
	wg := sync.WaitGroup{}
	wg.Add(len(*c.Descriptions))

	for i, _ := range *c.Descriptions {
		go executeCommand(c.Shell.Executable, c.Shell.Arguments, &(*c.Descriptions)[i], output, &wg)
	}

	go func() {
		for {
			select {
			case m := <-output:
				fmt.Println(m)
			}
		}
	}()

	wg.Wait()

	// Allow for any final messages to come through...
	time.Sleep(1)

	return nil
}

func executeCommand(shell string, shellArgs []string, description *config.GazeDesc, output chan string, wg *sync.WaitGroup) {
	cmd := exec.Command(shell, strings.Join(shellArgs, " "), description.Command)

	cmdReader, _ := cmd.StdoutPipe()
	cmdErrorReader, _ := cmd.StderrPipe()

	scanner := bufio.NewScanner(cmdReader)
	errScanner := bufio.NewScanner(cmdErrorReader)

	stdout := make(chan string, 1)
	stderr := make(chan string, 1)

	go func() {
		for {
			select {
			case value := <-stdout:
				output <- out.ColorText(fmt.Sprintf("[%s] %s", description.Name, value), description.ColorCode)
			case value := <-stderr:
				output <- out.ColorText(fmt.Sprintf("[%s | ERROR] %s", description.Name, value), description.ColorCode)
			}
		}

	}()

	go reader(scanner, stdout)
	go reader(errScanner, stderr)

	if err := cmd.Run(); err != nil {
		if description.AttemptedRetries < description.Retries {
			output <- out.ColorText(fmt.Sprintf("[gazer] Attempting retry for %s...", description.Name), description.ColorCode)
			description.AttemptedRetries += 1
			go executeCommand(shell, shellArgs, description, output, wg)
		} else {
			output <- out.ColorText(fmt.Sprintf("[gazer] Max retries for %s exceeded. Exiting...", description.Name), description.ColorCode)
			wg.Done()
		}
	} else {
		output <- out.ColorText(fmt.Sprintf("[gazer] Exiting %s...", description.Name), description.ColorCode)

		wg.Done()
	}

}

func reader(scanner *bufio.Scanner, out chan string) {
	for scanner.Scan() {
		out <- scanner.Text()
	}
}
