package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type result struct {
	url    string
	exists bool
}

func checkIfExists(ctx context.Context, urls <-chan string) <-chan result {
	responsec := make(chan result)

	go func() {
		defer close(responsec)

		for url := range urls {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				fmt.Printf("context error: %v\n", err)
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

func run(ctx context.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	urls := make(chan string, 4)
	// defer close(urls)

	urls <- "http://www.google.com"
	urls <- "http://www.facebook.com"
	urls <- "http://www.twitter.com"
	urls <- "http://in-valid-url.invalid"
	close(urls)

	c := checkIfExists(ctxWithTimeout, urls)

	now := time.Now()
	for result := range c {
		fmt.Printf("%s exists: %t\n", result.url, result.exists)
	}
	fmt.Printf("took %s\n", time.Since(now))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	run(ctx)
}
