package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10, // set connection pool size for each host
			MaxIdleConns:        100,
		},
	}
}

func startHTTPserver() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(50) * time.Microsecond)
		fmt.Fprintf(w, "Hello world")
	})

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

}

func startHTTPRequest(index int, wg *sync.WaitGroup) {
	counter := 0
	for i := 0; i < 10; i++ {
		resp, err := httpClient.Get("http://localhost:8080/")
		if err != nil {
			panic(fmt.Sprintf("Error: %v", err))
		}
		io.Copy(ioutil.Discard, resp.Body) // fully read the response body
		resp.Body.Close()                  // close the response body
		log.Printf("HTTP request #%v in Goroutine #%v", counter, index)
		counter += 1
		time.Sleep(time.Duration(1) * time.Second)
	}
	wg.Done()
}

func main() {
	startHTTPserver()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go startHTTPRequest(i, &wg)
	}
	wg.Wait()
}
