package main

import (
	"io"

	"golang.org/x/net/html"
)

// Job is the structure that defines the Job message object that is passed
// between goroutines. It describes a Job to be completed by one of them.
type Job struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ByteResult is the structure that defines the ByteResult message object that is passed
// between goroutines. It describes the Result of a job that is returned with the
// corresponding message in a byte array.
type ByteResult struct {
	Result []byte `json:"result"`
}

// JSONResult is the structure that defines the JSONResult message object that is passed
// between goroutines. It describes the Result of a job that is returned with the
// corresponding message in JSON notification (as a map[string]interface).
type JSONResult struct {
	Result map[string]interface{} `json:"result"`
}

// HTMLMeta is the structure that defines the HTMLMeta object and its properties
type HTMLMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"site_name"`
}

func extract(resp io.Reader) *HTMLMeta {
	z := html.NewTokenizer(resp)

	titleFound := false

	hm := new(HTMLMeta)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return hm
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == `body` {
				return hm
			}
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				desc, ok := extractMetaProperty(t, "description")
				if ok {
					hm.Description = desc
				}

				ogTitle, ok := extractMetaProperty(t, "og:title")
				if ok {
					hm.Title = ogTitle
				}

				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					hm.Description = ogDesc
				}

				ogImage, ok := extractMetaProperty(t, "og:image")
				if ok {
					hm.Image = ogImage
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					hm.SiteName = ogSiteName
				}
			}
		case html.TextToken:
			if titleFound {
				t := z.Token()
				hm.Title = t.Data
				titleFound = false
			}
		}
	}
	return hm
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}
