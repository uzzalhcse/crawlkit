package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Function to send a request and return the response code
func sendRequest(url string, index int, results chan<- string, successCounter *int, failureCounter *int, mutex *sync.Mutex) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error: %v", err)
		results <- fmt.Sprintf("Request %d to %s failed", index, url)
		mutex.Lock()
		*failureCounter++
		mutex.Unlock()
		return
	}
	defer resp.Body.Close()
	results <- fmt.Sprintf("Request %d to %s - Response Code: %d", index, url, resp.StatusCode)
	mutex.Lock()
	*successCounter++
	mutex.Unlock()
}

func main() {
	// Configurable parameters
	urls := []string{
		"https://markt-mall.jp/api/product/32",
		"https://markt-mall.jp/api/product/33",
		"https://markt-mall.jp/api/product/94",
		"https://markt-mall.jp/api/product/95",
		"https://markt-mall.jp/api/product/96",
		"https://markt-mall.jp/api/product/133",
		"https://markt-mall.jp/api/product/134",
		"https://markt-mall.jp/api/product/135",
		"https://markt-mall.jp/api/product/346",
	}
	startRequestRate := 7             // Start with 5 requests per second
	step := 1                         // Increment by 5 requests per second
	maxRate := 100                    // Upper limit of requests per second
	sleepDuration := 30 * time.Second // Sleep duration between tests

	requestRate := startRequestRate
	var totalSuccess, totalFailed int
	mutex := &sync.Mutex{}

	for requestRate <= maxRate {
		fmt.Printf("Testing with %d requests per second...\n", requestRate)

		var wg sync.WaitGroup
		results := make(chan string, requestRate*len(urls)) // Buffer size is rate * number of URLs

		startTime := time.Now()

		// Send requests at the current rate
		for _, url := range urls {
			for i := 0; i < requestRate; i++ {
				wg.Add(1)
				go func(url string, index int) {
					defer wg.Done()
					sendRequest(url, index, results, &totalSuccess, &totalFailed, mutex)
				}(url, i)
			}
		}

		// Wait for all goroutines to finish
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		for result := range results {
			fmt.Println(result)
		}

		elapsedTime := time.Since(startTime)
		fmt.Printf("Elapsed time: %v, Failed requests: %d, Successful requests: %d\n", elapsedTime, totalFailed, totalSuccess)

		// If there were failed requests, stop testing
		if totalFailed > 0 {
			fmt.Printf("Request limit reached at %d requests per second with %d failed requests.\n", requestRate, totalFailed)
			break
		}

		// Increase the request rate for the next iteration
		requestRate += step

		// Wait a bit before the next test to avoid overwhelming the server
		time.Sleep(sleepDuration)
	}
	fmt.Printf("Total Successful requests: %d\n", totalSuccess)
	fmt.Printf("Total Failed requests: %d\n", totalFailed)
}
