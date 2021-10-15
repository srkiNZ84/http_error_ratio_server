package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var rng = *rand.New(rand.NewSource(time.Now().UnixNano()))
var ratio int = 30
var slow bool = true
var slow_duration = 10
var errCounter int = 0
var totalCount int = 0
var shutdown_wait int = 20

const (
	success_message = "Everything's fine!"
	failure_message = "Server Error!"
	slow_message    = "Delayed response by %v seconds"
)

func main() {
	mux := http.NewServeMux()
	log.Println("Starting server...")
	mux.HandleFunc("/", randomHandler(slowHandler()))

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

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting to listen on port %s", port)
	log.Printf("Approximately %v percent of requests will return an HTTP 503", ratio)

	// Run our server in a go routine so that it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting http server %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	log.Println("Received our signal!!!")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(shutdown_wait))
	defer cancel()

	log.Println("Waiting for connections to finish up...")
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down %v", err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
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

func randomHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Do random response generation here
		res := returnRandomResponse(&rng, ratio)
		w.WriteHeader(res.status)
		fmt.Fprint(w, res.message)

		// call other function here
		fn(w, r)
	}
}

func slowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sr RandomResponse
		if slow {
			sr = returnSlowResponse(slow_duration)
		} else {
			sr = returnSlowResponse(0)
		}
		fmt.Fprintf(w, sr.message)
	}
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
