package main

import (
	"net/http"
	"testing"
	"time"
)

type NotRandomSource struct {
	n int
}

func (r NotRandomSource) Intn(i int) int {
	return r.n
}

func TestReturnSlowResponse(t *testing.T) {
	start := time.Now()
	returnSlowResponse(2)
	got := time.Since(start)
	want, _ := time.ParseDuration("1s")

	if got < want {
		t.Errorf("got a response in %v wanted more than %v", got, want)
	}
}

func TestReturnRandomResponse(t *testing.T) {
	cases := []struct {
		Description string
		FailRate    int
		RandomInt   int
		want        RandomResponse
	}{
		{"Return success when rand is clearly above", 30, 100, RandomResponse{http.StatusOK, success_message}},
		{"Return failure when rand is clearly below", 30, 1, RandomResponse{http.StatusServiceUnavailable, failure_message}},
		{"Return success from FailRate boundary", 30, 30, RandomResponse{http.StatusOK, success_message}},
		{"Return failure from FailRate boundary", 30, 29, RandomResponse{http.StatusServiceUnavailable, failure_message}},
		{"Return failure when rand is zero", 30, 0, RandomResponse{http.StatusServiceUnavailable, failure_message}},
		{"Always return success when error rate is zero", 0, 0, RandomResponse{http.StatusOK, success_message}},
		{"Always return success when error rate is zero", 0, 100, RandomResponse{http.StatusOK, success_message}},
		{"Always return success when error rate is zero", 0, 50, RandomResponse{http.StatusOK, success_message}},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := returnRandomResponse(NotRandomSource{test.RandomInt}, test.FailRate)
			if got != test.want {
				t.Errorf("got %q want %q", got, test.want)
			}
		})
	}
}
