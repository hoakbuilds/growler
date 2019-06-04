package main

import (
	"fmt"
	"strconv"
	"time"
)

const (
	unsplashSearchPhotosEndpoint = "https://unsplash.com/search/photos/"
)

func crawl(searchTerm string) error {

	b, err := requestURL(unsplashSearchPhotosEndpoint + searchTerm)

	if err != nil {
		return err
	}

	tm := time.Now().UnixNano() / 1000000
	filename := strconv.FormatInt(tm, 10) + ".html"

	err = saveRequest(b, filename)
	if err != nil {
		fmt.Printf("[error] there was an error saving the bytes to file %s.\n[GRWLR] > ", filename)
	}

	fmt.Printf("test:\n%v", b)

	return nil

}
