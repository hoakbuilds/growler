package main

// Job is the structure that defines the Job message object that is passed
// between goroutines. It describes a Job to be completed by one of them.
type Job struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
