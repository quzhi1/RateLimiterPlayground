package main

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	url         = "https://graph.microsoft.com/v1.0/me/messages?%24top=2&%24select=id"
	bearerToken = "Bearer ..." // Replace with your actual token
)

func main() {
	var wg sync.WaitGroup
	responseChan := make(chan string, 100)

	client := &http.Client{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(requestNumber int) {
			defer wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				responseChan <- fmt.Sprintf("Request %d: Failed to create request: %v", requestNumber, err)
				return
			}
			req.Header.Set("Authorization", bearerToken)

			resp, err := client.Do(req)
			if err != nil {
				responseChan <- fmt.Sprintf("Request %d: Failed to execute request: %v", requestNumber, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusTooManyRequests {
				retryAfter := resp.Header.Get("Retry-After")
				responseChan <- fmt.Sprintf("Request %d: 429 Too Many Requests. Retry-After: %s", requestNumber, retryAfter)
			} else {
				responseChan <- fmt.Sprintf("Request %d: Status %d.", requestNumber, resp.StatusCode)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(responseChan)
	}()

	// Collect and print results
	for res := range responseChan {
		fmt.Println(res)
	}
}
