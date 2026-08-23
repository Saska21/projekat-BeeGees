package main

import (
	"fmt"
	"net/http"
	"sync"
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

		for {
			select {
			case <-done:
				return
			case url, ok := <-urls:
				if !ok {
					return
				}
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

func merge[T any](done <-chan struct{}, cs ...<-chan T) <-chan T {
	results := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case r, ok := <-c:
					if !ok {
						return
					}
					results <- r
				}
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
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

	c1 := checkIfExists(done, urls)
	c2 := checkIfExists(done, urls)
	c3 := checkIfExists(done, urls)

	now := time.Now()
	for result := range merge(done, c1, c2, c3) {
		fmt.Printf("%s exists: %t\n", result.url, result.exists)
	}
	fmt.Printf("took %s\n", time.Since(now))

}
