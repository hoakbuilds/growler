package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	noFileOrDirectory string = " no such file or directory"
)

var (
	debug bool
	wg    sync.WaitGroup
)

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

	defer wg.Wait()

	var (
		tabCntl TabController
		wg      sync.WaitGroup
	)

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Logger started.")

	// before performing a normal start, check for a previous session
	restore, latestID, err := restoreSession()

	if err != nil {
		//assume there is no previous session
		errSplit := strings.Split(err.Error(), ":")
		if errSplit[1] == noFileOrDirectory {
			log.Printf("[restoreSession]: no previous session found")
		} else {
			log.Printf("[restoreSession]: %v\n", err)
		}
	}

	tabCntlChan := make(chan string, 10)
	shellCommChan := make(chan string, 10)

	if len(restore) != 0 {
		tabCntl = TabController{
			tcChan: tabCntlChan,
			shChan: shellCommChan,
			tabs:   []Tab{},
			id:     latestID,
		}

	} else {
		tabCntl = TabController{
			tcChan: tabCntlChan,
			shChan: shellCommChan,
			tabs:   []Tab{},
			id:     0,
		}
	}

	param := tabCntl.Run
	wg.Add(1)
	go param(restore)

	ReadForever(tabCntlChan, shellCommChan)
}

// ReadForever is a shell abstraction to get commands from the cli.
func ReadForever(tabCntlChan chan string, shellCommChan chan string) {
	log.Print("[readForever]\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("[growler] > ")
		// Read the keyboad input.
		in, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		input := processInput(in)

		if debug == true {
			log.Printf("[shell_debug]: %v\n", input)
		}

		if input[0] == "help" {
			fmt.Fprintf(os.Stdout, "\n \t Command usage help\n")
			fmt.Fprintf(os.Stdout, "\nget <url> \t - Opens a new tab and gives it the job of requesting the given url.\n")
			fmt.Fprintf(os.Stdout, "\nlist \t\t - Lists all open tabs, their current or last job, and their current status.\n")
			fmt.Fprintf(os.Stdout, "\nstop <tab id> \t - Stops the thread represented by the input id.\n\t\t If a thread is working it will wait until it has finished to stop it.\n")
			fmt.Fprintf(os.Stdout, "\njob <url> \t - Adds the given url to the pool of jobs, it will be executed once a tab is available.\n")
			fmt.Fprintf(os.Stdout, "\nexit \t\t - Stops all tabs and exits the browser.\n\n")
			continue
		} else if input[0] == "list" {

			tabCntlChan <- input[0]
			recv := <-shellCommChan

			if recv == "ok" {
				continue
			}

		} else if input[0] == "get" {

			if len(input) < 2 {
				fmt.Print("[error] you must provide a valid URL.\n")
				continue
			} else {
				tabCntlChan <- input[0]
				recv := <-shellCommChan

				if recv == "url" {
					tabCntlChan <- input[1]
				}

				ok := <-shellCommChan

				if ok == "ok" {
					continue
				}

			}

		} else if input[0] == "exit" {
			tabCntlChan <- input[0]
			msg := <-shellCommChan

			if msg == "ok" {
				break
			}

		}
		if input[0] == "stop" {

			if len(input) < 2 {
				fmt.Print("[error] you must provide a valid Tab ID.\n")
				continue
			} else {
				tabCntlChan <- input[0]
				recv := <-shellCommChan

				if recv == "TabID" {
					tabCntlChan <- input[1]
				}

				continue
			}

		}
		if input[0] == "stop" && input[1] == "all" {

			tabCntlChan <- input[0]
			recv := <-shellCommChan

			if recv == "TabID" {
				tabCntlChan <- input[1]
			}
			ok := <-shellCommChan
			if ok == "ok" {
				continue
			}

		}

		if input[0] == "job" {

			if err != nil {
				log.Printf("[err] %s", err)
			}

		}

	}

}
