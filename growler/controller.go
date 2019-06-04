package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	// GetInfo is the function used by the Controller interface to get information
	// about running tasks
	GetInfo()

	// NewTab is the function used by the Controller interface to spawn a new tab
	NewTab()

	// RelayJob is used to relay jobs to opened tabs that are not working
	RelayJob()

	// StopTab is used to stop a tab that is currently waiting for a job or
	// simply running. If the tab is running a job, it will be stopped once it
	// has finished that job.
	StopTab()

	// TerminateGracefully is the function used by the TabController
	// to issue a Stop request to all tabs and then terminate, allowing
	// the shell to exit gracefully.
	TerminateGracefully()

	// WriteSession is a function called by the browser to save current tab states
	WriteSession()

	// FinishRestoringSession is the function used by the TabController
	// to finish restoring a previous session
	FinishRestoringSession()
}

// TabController is the structure that defines the TabController
// and it's needed parameters
type TabController struct {
	// shChan is the channel used to send the response
	// to the shell
	shChan chan string
	// tcChan is the channel used to get the messages
	// that are sent from the shell
	tcChan chan string
	tabs   []Tab
	id     int
}

// Run is the function used by the TabController to start it's own lifecycle
func (c TabController) Run(restoredTabs []RestoredTab) {
	defer wg.Done()
	if len(restoredTabs) > 0 {
		newTabs, err := c.FinishRestoringSession(restoredTabs)
		if err != nil {
			log.Printf("%v", err)
		}
		c.tabs = newTabs
	}

	if c.id == 0 {
		c.id++
	}

	for {
		msg := <-c.tcChan

		if msg == "list" {

			log.Printf("[tabCntlDebug]: %v\n", msg)

			list := c.GetInfo()

			msg := strings.Join(list, "")

			fmt.Print(msg)
			c.shChan <- "ok"

		} else if msg == "get" {

			log.Printf("[tabCntlDebug]: %v\n", msg)
			c.shChan <- "url"

			jobMsg := <-c.tcChan

			log.Printf("[tabCntlDebug]: %v\n", jobMsg)
			job := Job{
				URL: jobMsg,
			}
			tab, err := c.NewTab(job)
			if err != nil {
				log.Printf("%v", err)
			}
			param := tab.Start
			wg.Add(1)
			go param()
			c.tabs = append(c.tabs, tab)
			c.id++

			c.shChan <- "ok"

		} else if msg == "stop" {
			var (
				newTabs []Tab
			)
			log.Printf("[tabCntlDebug]: %v\n", msg)

			c.shChan <- "TabID"
			msg := <-c.tcChan

			if msg == "all" {
				log.Printf("[tabCntlDebug]: %v\n", msg)
				for _, tab := range c.tabs {

					tabs, err := c.StopTab(tab.ID)
					if err != nil {
						log.Printf("[err] %s", err)
					}
					newTabs = tabs
				}

			} else {
				tabID, err := strconv.Atoi(msg)

				if err != nil {
					log.Printf("[err] %s", err)
				}

				log.Printf("[tabCntlDebug]: %v\n", tabID)

				tabs, err := c.StopTab(tabID)
				if err != nil {
					log.Printf("[err] %s", err)
				}
				newTabs = tabs

			}
			c.tabs = newTabs
			c.shChan <- "ok"

		} else if msg == "exit" {

			log.Printf("[tabCntlDebug]: saving session\n")
			err := c.WriteSession()

			if err != nil {
				panic(err)
			}

			for _, tab := range c.tabs {
				tab.jobsChan <- "stop"
				log.Printf("[tabCntlDebug]: notified tab %d to stop\n", tab.ID)
			}

			c.shChan <- "ok"
			break
		}

	}

	log.Printf("[tabCntlDebug] deferringDone\n")

}

// NewTab is the function used by the TabController to issue a new tab
// having it start it's lifecycle
func (c TabController) NewTab(job Job) (Tab, error) {

	var (
		tab    Tab
		jobURL string
	)

	j := make(chan string)
	r := make(chan string)

	if job.URL != "" {

		split := strings.Split(job.URL, "://")
		if split[0] != "https" {
			jobURL = "https://" + split[0]
		} else {
			jobURL = job.URL
		}

		tab = Tab{
			jobsChan:    j,
			resultsChan: r,
			ID:          c.id,
			Job:         jobURL,
			Status:      working,
		}

	}
	return tab, nil
}

// StopTab is the function used by the TabController to stop a running tab,
// given that it isn't currently working
func (c TabController) StopTab(ID int) ([]Tab, error) {

	log.Printf("[tabCntlDebug]: searching for tab to stop\n")
	i := 0
	for _, tab := range c.tabs {

		if ID == tab.ID {
			tab.jobsChan <- "stop"

			log.Printf("[tabCntlDebug]: stop message sent to %d\n", tab.ID)

			// Remove the element at index i from a.
			copy(c.tabs[i:], c.tabs[i+1:])  // Shift tabs[i+1:] left one index.
			c.tabs[len(c.tabs)-1] = Tab{}   // Erase last element (write zero value).
			c.tabs = c.tabs[:len(c.tabs)-1] // Truncate slice.

		}
		i++
	}
	return c.tabs, nil
}

// GetInfo is the function used by the TabController to give information about,
// running tabs and their jobs
func (c TabController) GetInfo() []string {
	var (
		lines []string
	)

	lines = append(lines, "\n Tab ID\tJob\t\t\tStatus\t\n")
	lines = append(lines, " ----\t----\t\t\t----\t\n")
	for _, tab := range c.tabs {
		string := fmt.Sprintf(" %d\t%s\t%s\t\n", tab.ID, tab.Job, tab.Status)
		lines = append(lines, string)
	}
	lines = append(lines, "\n")
	return lines
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

// FinishRestoringSession is the function used by the TabController
// to finish restoring a previous session
func (c TabController) FinishRestoringSession(restoredTabs []RestoredTab) ([]Tab, error) {
	var (
		newTabs []Tab
	)
	for _, tab := range restoredTabs {
		job := Job{
			URL: tab.Job,
		}
		newTab, err := c.NewTab(job)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}
		newTab.Status = waitingJob
		param := newTab.Start
		wg.Add(1)
		go param()
		newTabs = append(newTabs, newTab)
		c.id++
	}

	return newTabs, nil
}

// WriteSession is a function called by the browser to save current tab states
func (c TabController) WriteSession() error {
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
