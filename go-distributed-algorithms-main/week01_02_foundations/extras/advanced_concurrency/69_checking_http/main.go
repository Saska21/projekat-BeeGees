package main

import (
	"fmt"
	"net/http"
	"time"
)

type result struct {
	url    string
	exists bool
}

func checkIfExists(done <-chan struct{}, urls <-chan string) <-chan result {
	responsec := make(chan result)

	go func() {
		defer close(responsec)

		for url := range urls {
			select {
			case <-done:
				return
			default:
				res, err := http.Get(url)
				if err != nil {
					responsec <- result{url: url, exists: false}
				} else if res.StatusCode == http.StatusOK {
					responsec <- result{url: url, exists: true}
				} else {
					responsec <- result{url: url, exists: false}
				}
			}
		}
	}()

	return responsec
}

func main() {
	done := make(chan struct{})
	defer close(done)

	urls := make(chan string, 4)
	// defer close(urls)

	urls <- "http://www.google.com"
	urls <- "http://www.facebook.com"
	urls <- "http://www.twitter.com"
	urls <- "http://in-valid-url.invalid"
	close(urls)

	c := checkIfExists(done, urls)

	now := time.Now()
	for result := range c {
		fmt.Printf("%s exists: %t\n", result.url, result.exists)
	}
	fmt.Printf("took %s\n", time.Since(now))

}
