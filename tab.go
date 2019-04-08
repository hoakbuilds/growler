package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type tabExecutionStatus int

const (
	//
	waitingJob tabExecutionStatus = 1
	working    tabExecutionStatus = 2
)

// Tab is the structure that defines a Tab goroutine and all its needed params
type Tab struct {
	jobsChan    chan string
	resultsChan chan string
	ID          int
	Job         string
	Status      tabExecutionStatus
}

// requestURL receives a url in the form of a string and returns
// a []byte with the byte content of that request's
// response
func requestURL(url string) ([]byte, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// saveRequest saves the byte content of a request's response to a file
func saveRequest(content []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// write the content to the file
	if _, err = f.Write(content); err != nil {
		return err
	}
	return nil
}

func (status tabExecutionStatus) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)

	strings := []string{
		waitingJob: "Waiting for job.",
		working:    "Working..",
	}

	// return the string constant
	// from the status array above.
	return strings[status]
}

// restoreSession
func restoreSession() ([]Tab, error) {
	var (
		tabs []Tab
	)

	data, err := ioutil.ReadFile("tab.data")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), ",")
	for _, line := range lines {

		fmt.Print(line)
	}
	//
	f, err := os.OpenFile("tab.data", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	// nuke the previous session
	if _, err = f.WriteString(" "); err != nil {
		panic(err)
	}
	return tabs, nil
}

// writeFile is a function called by the browser to save current tab states
func writeFile(tb []Tab) error {
	f, err := os.OpenFile("tab.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	for _, tab := range tabs {

		//write tab id
		if _, err = f.WriteString("Tab ID: "); err != nil {
			panic(err)
		}
		id := strconv.Itoa(tab.ID)
		if _, err = f.WriteString(id); err != nil {
			panic(err)
		}

		//write job
		if _, err = f.WriteString(",\nJob: "); err != nil {
			panic(err)
		}
		if _, err = f.WriteString(tab.Job); err != nil {
			panic(err)
		}

		//write status
		if _, err = f.WriteString(",\nStatus: "); err != nil {
			panic(err)
		}
		var status string
		if tab.Status == 1 {
			status = "Waiting Job"
		}
		if tab.Status == 2 {
			status = "Working"
		}
		if _, err = f.WriteString(status); err != nil {
			panic(err)
		}

		/*write last active session
		t := time.Now()
		if _, err = f.WriteString(",\nLast time active: "); err != nil {
			panic(err)
		}

		if _, err = f.WriteString(t.Format("2006-01-02 15:04:05")); err != nil {
			panic(err)
		}*/

		if _, err = f.WriteString(",\n-------------------- || --------------------\n"); err != nil {
			panic(err)
		}

		if err != nil {
			fmt.Println(err)
			f.Close()
		}

	}

	return nil
}

// Start is the starting function for a jobless tab.
func (t Tab) Start() {
	defer wg.Done()
	log.Printf("[tab-%d-start]\n", t.ID)

	for {
		log.Printf("[tab-%d] starting %s\n", t.ID, t.Job)
		for {

			elapsed := time.Now().UnixNano() / 1000000

			b, err := requestURL(t.Job)

			if err != nil {
				fmt.Printf("[error] there was an error requesting %s.\n[growler] > ", t.Job)
			} else {

				tm := time.Now().UnixNano() / 1000000
				filename := strconv.FormatInt(tm, 10) + ".html"

				err := saveRequest(b, filename)
				if err != nil {
					fmt.Printf("[error] there was an error saving the bytes to file %s.\n[growler] > ", filename)
				} else {
					fmt.Printf("\n[success] request to %s saved to file %s\nTime elapsed: %d ms\nRequest size: %d bytes\n[growler] > ", t.Job, filename, (tm - elapsed), len(b))
				}

				break
			}

		}

		t.Status = waitingJob
		log.Printf("[tab-%d-waitingJob]\n", t.ID)
		job := <-t.jobsChan

		if job == "stop" {
			log.Printf("[tab-%d-stopping]\n", t.ID)
			break
		}

		if t.Job != job {
			log.Printf("[tab-%d-setting] job: %s, new job: %s\n", t.ID, t.Job, job)
			t.Job = job
			t.Status = working
		}

		log.Printf("[tab-%d] end - %s\n", t.ID, job)

	}
	log.Printf("[tab-%d-deferringDone]\n", t.ID)
}
