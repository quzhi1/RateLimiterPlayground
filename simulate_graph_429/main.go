package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

const (
	url         = "https://graph.microsoft.com/v1.0/me/messages?%24top=2&%24select=id"
	bearerToken = "Bearer ..." // Replace with your actual token
)

const (
	numRequests = 20
)

func main() {
	var wg sync.WaitGroup
	responseChan := make(chan string, numRequests)

	client := &http.Client{}

	for i := 0; i < numRequests; i++ {
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
				headers := ""
				for k, v := range resp.Header {
					headers += fmt.Sprintf("%s: %v\n", k, v)
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					responseChan <- fmt.Sprintf("Request %d: Failed to read response body: %v", requestNumber, err)
					return
				}
				responseStr := string(body)
				responseChan <- fmt.Sprintf("Request %d: 429 Too Many Requests. Headers: %s\nResponse: %s", requestNumber, headers, responseStr)
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
