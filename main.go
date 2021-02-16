package main

import (
	"fmt"
	"math/rand"
	"time"
)

var ration int = 30

func main() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	rv := []int{}
	var downCounter int = 0
	for i := 0; i < 100; i++ {
		t := r.Intn(100)
		rv = append(rv, t)
		if t >= ration {
			fmt.Println("200 OK")
		} else {
			downCounter++
			fmt.Println("503 Service Unavailable")

		}
	}

	fmt.Println("Here it is", rv)
	fmt.Println("Downcounter is", downCounter)
}
