package main

import (
	"fmt"
	"os"
	"strconv"
)

// Job is the structure that defines the Job message object that is passed
// between goroutines. It describes a Job to be completed by one of them.
type Job struct {
	URL string `json:"url"`
}

// Controller interface, which the TabController implements
type Controller interface {
	// Run is the function used by the Controller interface to execute its
	// tasks, which are to get commands from the shell and relay their
	// consequent jobs
	Run()
	NewTab()
	RelayJob()
	StopTab()
	TerminateGracefully()
	WriteSession()
}

// TabController is the structure that defines the TabController
// and it's needed parameters
type TabController struct {
	// shChan is the channel used to get the messages
	// that are sent from the shell
	shChan chan string
	tabs   []Tab
	id     int
}

// Run is the function used by the TabController to start it's own lifecycle
func (c TabController) Run() {

	c.id++

	for {
		msg := <-c.shChan
		fmt.Print(msg)

	}

}

// NewTab is the function used by the TabController to issue a new tab
// having it start it's lifecycle
func (c TabController) NewTab() ([]Tab, error) {
	return c.tabs, nil
}

// StopTab is the function used by the TabController to stop a running tab,
// given that it isn't currently working
func (c TabController) StopTab() ([]Tab, error) {
	return c.tabs, nil
}

// RelayJob is the function used by the TabController look up a tab that is not
// currently working and to relay a job to it.
func (c TabController) RelayJob(job Job) error {

	for _, tab := range c.tabs {
		if tab.Status == waitingJob {
			tab.jobsChan <- string(job.URL)
		}
	}

	return nil
}

// TerminateGracefully is the function used by the TabController
// to issue a Stop request to all tabs and then terminate, allowing
// the shell to exit gracefully.
func (c TabController) TerminateGracefully() error {
	return nil
}

// WriteSession is a function called by the browser to save current tab states
func (c TabController) WriteSession(tb []Tab) error {
	f, err := os.OpenFile("tab.data", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	for _, tab := range c.tabs {

		//write tab id
		if _, err = f.WriteString("Tab ID::"); err != nil {
			panic(err)
		}
		id := strconv.Itoa(tab.ID)
		if _, err = f.WriteString(id); err != nil {
			panic(err)
		}

		//write job
		if _, err = f.WriteString(",\nJob::"); err != nil {
			panic(err)
		}
		if _, err = f.WriteString(tab.Job); err != nil {
			panic(err)
		}

		//write status
		if _, err = f.WriteString(",\nStatus::"); err != nil {
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

		if _, err = f.WriteString("\n--------------------------------------------\n"); err != nil {
			panic(err)
		}

		if err != nil {
			fmt.Println(err)
			f.Close()
		}

	}

	return nil
}
