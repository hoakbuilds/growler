package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	debug bool
)

type tabExecutionStatus int

const (
	//
	startingExecution tabExecutionStatus = 0
	waitingJob        tabExecutionStatus = 1
	working           tabExecutionStatus = 2
	done              tabExecutionStatus = 3
)

// Tab is the structure that defines a Tab goroutine and all its needed params
type Tab struct {
	jobsChan    chan string
	resultsChan chan string
	ID          int
	job         string
	Status      tabExecutionStatus
}

// SetJob sets the objects job field
func (t *Tab) SetJob(job string) Tab {
	t.job = job
	return *t
}

// Job gets the objects job field
func (t Tab) Job() string {
	return t.job
}

func (status tabExecutionStatus) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)

	strings := []string{
		startingExecution: "Starting execution.",
		waitingJob:        "Waiting for job.",
		working:           "Working..",
		done:              "Done",
	}

	// return the string constant
	// from the status array above.
	return strings[status]
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

	var wg sync.WaitGroup

	ReadForever(&wg)

	return
}

// ReadForever is a shell abstraction to get commands from the cli.
func ReadForever(wg *sync.WaitGroup) {
	log.Print("[readForever]\n")

	var tabs []Tab
	var i = 0
	reader := bufio.NewReader(os.Stdin)

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
			j := make(chan string)
			r := make(chan string)
			tab := Tab{
				jobsChan:    j,
				resultsChan: r,
				ID:          i,
				job:         "n/a",
				Status:      startingExecution,
			}

			wg.Add(1)
			go tab.Start(wg)
			tabs = append(tabs, tab)
			i++

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
			for _, tab := range tabs {

				id, err := strconv.Atoi(input[1])

				if err != nil {
					log.Panicf("[err] %s", err)
				}

				if id == tab.ID {
					tab.jobsChan <- input[2]

				}

			}

		}

		if input[0] == "exit" {
			break
		}

	}

}

// Start is the starting function for a jobless tab.
func (t Tab) Start(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Printf("[tab-%d-start]\n", t.ID)

	job := <-t.jobsChan

	if t.Job() != job {
		(&t).SetJob(job)
	}

	log.Printf("[tab-%d] %s\n", t.ID, job)

	for {
		log.Printf("[tab-%d] start - %s\n", t.ID, job)

		time.Sleep(2 * time.Second)
		break

	}

	log.Printf("[tab-%d] end - %s\n", t.ID, job)

	(&t).SetJob("none....")

	return nil
}
