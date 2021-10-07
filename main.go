package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var rng = *rand.New(rand.NewSource(time.Now().UnixNano()))
var ratio int = 30
var errCounter int = 0
var totalCount int = 0

const (
	success_message = "Everything's fine!"
	failure_message = "Server Error!"
)

func main() {

	log.Println("Starting server...")
	http.HandleFunc("/", handler)

	ratio, _ := strconv.Atoi(os.Getenv("ERROR_RATIO"))
	if ratio == 0 {
		ratio = 30
		log.Printf("Default ratio of %v percent will be used", ratio)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Configuring default port %s", port)
	}

	log.Printf("Starting to listen on port %s", port)
	log.Printf("Approximately %v percent of requests will return an HTTP 503", ratio)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting http server", err)
	}
}

// Interface wrapper around things that can return a random integer
// Implemented by rand and NotRandomSource in the tests
type EntropySource interface {
	Intn(int) int
}

// RandomResponse type to capture the HTTP response returned from the server
type RandomResponse struct {
	status  int
	message string
}

func handler(w http.ResponseWriter, r *http.Request) {
	res := returnRandomResponse(&rng, ratio)
	w.WriteHeader(res.status)
	fmt.Fprint(w, res.message)
}

func returnRandomResponse(e EntropySource, r int) RandomResponse {
	num := e.Intn(100)
	if num >= r {
		totalCount++
		return RandomResponse{status: http.StatusOK, message: success_message}
	}

	errCounter++
	totalCount++
	return RandomResponse{status: http.StatusServiceUnavailable, message: failure_message}
}

func printStats(w http.ResponseWriter) {
	fmt.Fprintf(w, "\n\nTotal requests %v", totalCount)
	fmt.Fprintf(w, "\nTotal errors returned so far is %v", errCounter)
	fmt.Fprintf(w, "\nPercent errors returned so far is %v", float32(errCounter)/float32(totalCount))
}
