package main

import (
	"fmt"
    "net/http"
    "os"
    "strconv"
    "sync"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run test_script_rate_limiter.go <number_of_requests>")
        return
    }

    numRequests, err := strconv.Atoi(os.Args[1])
    if err != nil || numRequests <= 0 {
        fmt.Println("Please provide a valid positive integer for the number of requests.")
        return
    }
	url := "http://localhost:8080/ping"
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Request %d: Error: %v\n", i+1, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Request %d: Status: %d\n", i+1, resp.StatusCode)
		}(i)
	}

	wg.Wait()
}
