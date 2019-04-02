package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/imroc/req"
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
func (t *Tab) SetJob(job string) {
	t.job = job
}

// Job gets the objects job field
func (t Tab) Job() string {
	return t.job
}

// toReadableDate receives an interface with a timestamp
// and returns a time.Time structure, which is human-readable
func toReadableDate(timestamp interface{}) time.Time {

	sec, dec := math.Modf(timestamp.(float64))
	tm := time.Unix(int64(sec), int64(dec*(1e9)))

	return tm
}

func requestURL(url string) []byte {
	var contents []byte
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%s\n", string(contents))

	}
	return contents
}

// requestURL receives a url in the form of a string and returns
// a map[string]interface{} with the JSON content of that request's
// response
func requestAPI(url string) map[string]interface{} {
	// use Req object to initiate requests.
	req := req.New()
	req.Get(url)

	// use req package to initiate request.
	r, err := req.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	r.ToJSON(&res) // response => struct/map

	return res

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

// Start is the starting function for a jobless tab.
func (t Tab) Start() {
	defer wg.Done()
	log.Printf("[tab-%d-start]\n", t.ID)

	for {

		t.Status = waitingJob
		log.Printf("[tab-%d-waitingJob]\n", t.ID)
		job := <-t.jobsChan

		if job == "stop" {
			log.Printf("[tab-%d-stopping]\n", t.ID)
			break
		}

		if t.Job() != job {
			log.Printf("[tab-%d-setting] job: %s, new job: %s\n", t.ID, t.Job(), job)
			t.SetJob(job)
			t.Status = working
		}

		log.Printf("[tab-%d] starting %s\n", t.ID, job)

		for {
			log.Printf("[tab-%d] start - %s\n", t.ID, job)

			time.Sleep(2 * time.Second)
			break

		}

		t.Status = done

		log.Printf("[tab-%d] end - %s\n", t.ID, job)

		t.SetJob("none....")

	}
	log.Printf("[tab-%d-deferringDone]\n", t.ID)
}
