package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	debug bool
	tabs  []Tab
	id    int
	wg    sync.WaitGroup
)

//TabController is
type TabController interface {
	Start()
}

// checkArguments checks the CLI arguments passed and returns them as an array
func checkArguments() (arguments []string) {

	arguments = os.Args

	if len(arguments) != 1 {

		if arguments[1] == "debug" {
			debug = true
		}

		return arguments

	}

	return
}

// processInput will receive a command to execute and will perform logic prior to
// executing those commands, so that they are properly executed and terminated
// when it's needed to
func processInput(input string) []string {

	log.Print("[processInput]\n")
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	inputArr := strings.Split(input, " ")

	return inputArr
}

func init() {

	checkArguments()

}

// main is the main function of the program
func main() {

	if debug == true {
		f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
		log.Println("Logger started.")
	}

	ReadForever()

	wg.Wait()

	return
}

// ReadForever is a shell abstraction to get commands from the cli.
func ReadForever() {
	log.Print("[readForever]\n")

	reader := bufio.NewReader(os.Stdin)
	id++

	for {
		fmt.Print("[crawler] > ")
		// Read the keyboad input.
		in, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		input := processInput(in)

		if debug == true {
			log.Printf("[shell_debug]: %v\n", input)
		}

		if input[0] == "get" {

			var (
				job string
				tab Tab
			)

			if len(input) >= 2 {
				job = input[1]
			}

			j := make(chan string)
			r := make(chan string)

			if job != "" {

				tab = Tab{
					jobsChan:    j,
					resultsChan: r,
					ID:          id,
					job:         job,
					Status:      startingExecution,
				}

			} else {
				tab = Tab{
					jobsChan:    j,
					resultsChan: r,
					ID:          id,
					job:         "n/a",
					Status:      startingExecution,
				}
			}
			param := tab.Start
			wg.Add(1)
			go param()
			tabs = append(tabs, tab)
			id++

		}

		if input[0] == "list" {
			fmt.Fprintf(os.Stdout, "\n %s\t%s\t%s\t", "Tab ID", "Job", "Status")
			fmt.Fprintf(os.Stdout, "\n %s\t%s\t%s\t", "----", "----", "----")
			for _, tab := range tabs {
				fmt.Fprintf(os.Stdout, "\n %d\t%s\t%s\t", tab.ID, tab.Job(), tab.Status)
			}
			fmt.Fprintf(os.Stdout, "\n\n")

		}

		if input[0] == "assign" {
			id, err := strconv.Atoi(input[1])

			if err != nil {
				log.Panicf("[err] %s", err)
			}

			for _, tab := range tabs {

				if id == tab.ID {
					tab.jobsChan <- input[2]
					tab.SetJob(input[2])
				}

			}

		}

		if input[0] == "exit" {
			break
		}

		if (input[0] == "stop") && (len(input) == 2) {
			id, err := strconv.Atoi(input[1])

			if err != nil {
				log.Panicf("[err] %s", err)
			}
			i := 0
			for _, tab := range tabs {

				if id == tab.ID {
					tab.jobsChan <- "stop"

					// Remove the element at index i from a.
					copy(tabs[i:], tabs[i+1:]) // Shift tabs[i+1:] left one index.
					tabs[len(tabs)-1] = Tab{}  // Erase last element (write zero value).
					tabs = tabs[:len(tabs)-1]  // Truncate slice.

				}
				i++
			}
		}

	}

}
