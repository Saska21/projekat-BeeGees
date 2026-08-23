package main

import (
	"fmt"
	"net/http"
)

func checkIfExist(done <-chan struct{}, urls ...string) (<-chan *http.Response, <-chan error) {
	respChan := make(chan *http.Response)
	errChan := make(chan error)

	go func() {
		for _, url := range urls {
			select {
			case <-done:
				return
			default:
				resp, err := http.Get(url)
				if err != nil {
					errChan <- err
				} else {
					respChan <- resp
				}
			}
		}
		close(respChan)
		close(errChan)
	}()
	return respChan, errChan
}

func main() {
	done := make(chan struct{})

	respChan, errChan := checkIfExist(done, "https://www.google.com", "http://localhost:300")

	for range 2 {
		select {
		case res := <-respChan:
			fmt.Println("Status: ", res.Status)
		case err := <-errChan:
			fmt.Println("Error: ", err)
		}
	}

	close(done)
}
