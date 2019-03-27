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

type tab struct {
	jobsChan    chan string
	resultsChan chan string
	id          int
	job         string
	status      tabExecutionStatus
}

func (t *tab) setJob(job string) error {

	t.job = job

	return nil

}

func (t *tab) getJob() string {
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

	// return the name of a Weekday
	// constant from the names array
	// above.
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

	shellChan := make(chan string)

	go readForever(shellChan, &wg)

	for {
		a := <-shellChan

		if a == "exit" {
			break
		}
	}

	wg.Wait()

	return
}

func readForever(msg chan string, wg *sync.WaitGroup) {
	log.Print("[readForever]\n")

	reader := bufio.NewReader(os.Stdin)
	tabs := []tab{}
	i := 1

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
			tab := &tab{
				jobsChan:    j,
				resultsChan: r,
				id:          i,
				job:         "none",
				status:      startingExecution,
			}

			wg.Add(1)
			go tab.start(wg)
			tabs = append(tabs, *tab)
			i++

		}

		if input[0] == "list" {
			fmt.Fprintf(os.Stdout, "\n %s\t%s\t%s\t", "Tab ID", "Job", "Status")
			fmt.Fprintf(os.Stdout, "\n %s\t%s\t%s\t", "----", "----", "----")
			for _, tab := range tabs {
				fmt.Fprintf(os.Stdout, "\n %d\t%s\t%s\t", tab.id, tab.getJob(), tab.status)
			}
			fmt.Fprintf(os.Stdout, "\n")

		}

		if input[0] == "assign" {
			for _, tab := range tabs {

				id, err := strconv.Atoi(input[1])

				if err != nil {
					log.Panicf("[err] %s", err)
				}

				if id == tab.id {
					tab.jobsChan <- input[2]
				}

			}

		}

		msg <- in

	}

}

func (t *tab) start(wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Printf("[tab-%d-start]\n", t.id)

	job := <-t.jobsChan

	log.Printf("[tab-%d-job] %s\n", t.id, job)

	t.setJob(job)

	for {
		log.Printf("[tab-%d-job] start - %s\n", t.id, job)

		time.Sleep(2 * time.Second)
		break

	}

	return nil
}
