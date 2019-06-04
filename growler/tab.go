package main

import (
	"fmt"
	"io"
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

// RestoredTab is the structure that defines a Tab goroutine and all its needed params
type RestoredTab struct {
	jobsChan    chan string
	resultsChan chan string
	ID          int
	Job         string
	Status      tabExecutionStatus
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
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
func restoreSession() ([]RestoredTab, int, error) {
	var (
		tabs     []RestoredTab
		latestID int
	)

	data, err := ioutil.ReadFile("tab.data")
	if err != nil {
		return nil, 0, err
	}

	lines := strings.Split(string(data), "\n--------------------------------------------\n")

	for _, line := range lines {
		var (
			id  int
			job string
		)

		split := strings.Split(line, ",")
		if len(split) > 1 {
			for i := 0; i < len(split); i++ {

				params := strings.Split(split[i], "::")
				if i == 0 {
					id, err = strconv.Atoi(params[1])
					if id > latestID {
						latestID = id
					}

					if err != nil {
						log.Printf("[err] %s", err)
					}
				}
				if i == 1 {
					job = params[1]
				}
			}

			t := RestoredTab{
				ID:     id,
				Job:    job,
				Status: waitingJob,
			}
			tabs = append(tabs, t)
		}

	}
	return tabs, latestID, nil
}

// Start is the starting function for a jobless tab.
func (t Tab) Start() {
	defer wg.Done()
	log.Printf("[tab-%d-start]\n", t.ID)

	for {
		log.Printf("[tab-%d] starting %s\n", t.ID, t.Job)
		if t.Status != waitingJob && t.Job != "" {
			for {

				elapsed := time.Now().UnixNano() / 1000000

				b, err := requestURL(t.Job)

				if err != nil {
					fmt.Printf("[error] there was an error requesting %s.\n[GRWLR] > ", t.Job)
					break
				} else {

					tm := time.Now().UnixNano() / 1000000
					filename := strconv.FormatInt(tm, 10) + ".html"

					err := saveRequest(b, filename)
					if err != nil {
						fmt.Printf("[error] there was an error saving the bytes to file %s.\n[GRWLR] > ", filename)
					} else {
						fmt.Printf("\n[success] request to %s saved to file %s\nTime elapsed: %d ms\nRequest size: %d bytes\n[GRWLR] > ", t.Job, filename, (tm - elapsed), len(b))
					}

					break
				}

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
