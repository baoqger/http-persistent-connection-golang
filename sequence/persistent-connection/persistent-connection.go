package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func startHTTPServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(50) * time.Microsecond)
		fmt.Fprintf(w, "Hello world")
	})

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

}

func startHTTPRequest() {
	counter := 0
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			panic(fmt.Sprintf("Error: %v", err))
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		log.Printf("HTTP request #%v", counter)
		counter += 1
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func main() {
	startHTTPServer()

	startHTTPRequest()
}
