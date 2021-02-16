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

var rng rand.Rand
var ratio int = 30
var errCounter int = 0
var totalCount int = 0

func main() {
	log.Println("Seeding random number generator...")
	source := rand.NewSource(time.Now().UnixNano())
	rng = *rand.New(source)

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
	log.Printf("%v percent of requests will return an HTTP 503", ratio)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting http server", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	num := rng.Intn(100)
	if num >= ratio {
		w.WriteHeader(http.StatusOK)
		totalCount++
		fmt.Fprint(w, "Everything's fine!")
		printStats(w)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		errCounter++
		totalCount++
		fmt.Fprint(w, "Server Error!")
		printStats(w)
	}
}

func printStats(w http.ResponseWriter) {
	fmt.Fprintf(w, "\n\nTotal requests %v", totalCount)
	fmt.Fprintf(w, "\nTotal errors returned so far is %v", errCounter)
	fmt.Fprintf(w, "\nPercent errors returned so far is %v", float32(errCounter)/float32(totalCount))
}
