package gaze

import (
	"bufio"
	"fmt"
	"github.com/jettdc/gazer/config"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"
)

func ExecuteDescriptions(shellConfig *config.ShellConfig, descriptions *[]config.GazeDesc) error {
	output := make(chan string, 8)
	wg := sync.WaitGroup{}
	wg.Add(len(*descriptions))

	for i, _ := range *descriptions {
		go executeCommand(shellConfig.Executable, shellConfig.Arguments, &(*descriptions)[i], output, &wg)
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
				output <- description.Output(value)
			case value := <-stderr:
				output <- description.ErrorOutput(value)
			case <-description.Context.Done():
				return
			}
		}

	}()

	// Pass process output to the channels
	go reader(scanner, stdout)
	go reader(errScanner, stderr)

	cmd.Start()

	c := make(chan os.Signal, 1)

	manuallyKilled := false

	go func() {
		for {
			select {
			case s := <-c:
				manuallyKilled = true
				cmd.Process.Signal(s)
			case <-description.Context.Done():
				output <- description.AdminOutput(fmt.Sprintf("Cancel requested, killing %s", description.Name))
				cmd.Process.Signal(os.Kill)
				manuallyKilled = true
				return
			}
		}
	}()

	// Pass along signals to cmd
	signal.Notify(c, os.Interrupt, os.Kill)

	if err := cmd.Wait(); err != nil {
		if manuallyKilled { // Don't try to restart
			output <- description.AdminOutput(fmt.Sprintf("Process %s manually killed.", description.Name))
			wg.Done()
		} else if description.Restart == "always" || (description.Restart == "retry" && description.AttemptedRetries < description.Retries) {
			output <- description.AdminOutput(fmt.Sprintf("Attempting retry for %s...", description.Name))
			description.AttemptedRetries += 1
			go executeCommand(shell, shellArgs, description, output, wg)
		} else {
			if description.Restart != "" {
				output <- description.AdminOutput(fmt.Sprintf("Max retries for %s exceeded. Exiting...", description.Name))
			}
			wg.Done()
		}
	} else {
		output <- description.AdminOutput(fmt.Sprintf("Completed execution of %s", description.Name))
		description.AttemptedRetries = 0 // Reset in case it crashes again
		wg.Done()
	}
}

func reader(scanner *bufio.Scanner, out chan string) {
	for scanner.Scan() {
		out <- scanner.Text()
	}
}
