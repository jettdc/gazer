package gaze

import (
	"bufio"
	"fmt"
	"github.com/jettdc/gazer/config"
	"github.com/jettdc/gazer/out"
	"os/exec"
	"strings"
	"sync"
)

func ExecuteDescriptions(c config.Config) error {
	output := make(chan string, 8)
	wg := sync.WaitGroup{}
	wg.Add(len(*c.Descriptions))

	for _, desc := range *c.Descriptions {
		go executeCommand(c.Shell.Executable, c.Shell.Arguments, desc, output, &wg)
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
	return nil
}

func executeCommand(shell string, shellArgs []string, description config.GazeDesc, output chan string, wg *sync.WaitGroup) {
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
				output <- out.ColorText(fmt.Sprintf("[%s] [ERROR] %s", description.Name, value), description.ColorCode)
			}
		}

	}()

	go reader(scanner, stdout)
	go reader(errScanner, stderr)

	cmd.Run()

	output <- out.ColorText(fmt.Sprintf("[%s] Exiting...", description.Name), description.ColorCode)

	wg.Done()
}

func reader(scanner *bufio.Scanner, out chan string) {
	for scanner.Scan() {
		out <- scanner.Text()
	}
}
