package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/justinas/alice"
)

var rng = *rand.New(rand.NewSource(time.Now().UnixNano()))
var ratio int = 30
var slow bool = true
var slow_duration = 2
var errCounter int = 0
var totalCount int = 0
var port string = "8080"

const (
	success_message = "Everything's fine!"
	failure_message = "Server Error!"
	slow_message    = "Delayed response by %v seconds"
)

func main() {
	configure()

	chain := alice.New(randomHandler, slowHandler).Then(http.HandlerFunc(okHandler))

	log.Println("Starting server...")
	log.Printf("Starting to listen on port %s", port)
	log.Printf("Approximately %v percent of requests will return an HTTP 503", ratio)
	if err := http.ListenAndServe(":"+port, chain); err != nil {
		log.Fatal("Error starting http server", err)
	}
}

func configure() {
	ratio, _ := strconv.Atoi(os.Getenv("ERROR_RATIO"))
	if ratio == 0 {
		ratio = 30
		log.Printf("Default ratio of %v percent will be used", ratio)
	}

	slow_responses := os.Getenv("SLOW_RESPONES")
	if slow_responses != "" {
		slow = true
		log.Println("Slow responses are turned on")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Configuring default port %s", port)
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

func okHandler(w http.ResponseWriter, r *http.Request) {
	printStats(w)
}

func randomHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do random response generation here
		res := returnRandomResponse(&rng, ratio)
		w.WriteHeader(res.status)
		fmt.Fprint(w, res.message)
	})
}

func slowHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sr RandomResponse
		if slow {
			sr = returnSlowResponse(slow_duration)
		} else {
			sr = returnSlowResponse(0)
		}
		fmt.Fprintf(w, sr.message)
	})
}

func returnRandomResponse(e EntropySource, r int) RandomResponse {
	log.Printf("returning random response")
	num := e.Intn(100)
	if num >= r {
		totalCount++
		return RandomResponse{status: http.StatusOK, message: success_message}
	}

	errCounter++
	totalCount++
	return RandomResponse{status: http.StatusServiceUnavailable, message: failure_message}
}

func returnSlowResponse(s int) RandomResponse {
	log.Printf("In the slow response function")
	time.Sleep(time.Second * time.Duration(s))
	return RandomResponse{status: http.StatusOK, message: fmt.Sprintf(slow_message, s)}
}

func printStats(w http.ResponseWriter) {
	fmt.Fprintf(w, "\n\nTotal requests %v", totalCount)
	fmt.Fprintf(w, "\nTotal errors returned so far is %v", errCounter)
	fmt.Fprintf(w, "\nPercent errors returned so far is %v", float32(errCounter)/float32(totalCount))
}
