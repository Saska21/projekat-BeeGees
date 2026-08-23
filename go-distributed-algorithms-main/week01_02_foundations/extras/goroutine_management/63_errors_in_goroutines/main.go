package main

import (
	"fmt"
	"net/http"
)

type result struct {
	resp *http.Response
	err  error
}

func checkIfExist(done <-chan struct{}, urls ...string) <-chan result {
	results := make(chan result)

	go func() {
		for _, url := range urls {
			select {
			case <-done:
				return
			default:
				resp, err := http.Get(url)
				results <- result{resp: resp, err: err}
			}
		}
		close(results)
	}()
	return results
}

func main() {
	done := make(chan struct{})

	results := checkIfExist(done, "https://www.google.com", "http://localhost:300")

	for r := range results {
		if r.resp != nil {
			fmt.Println(r.resp.Status)
		}
		if r.err != nil {
			fmt.Println("Error:", r.err)
		}
	}

	close(done)
}
