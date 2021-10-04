package main

import (
	"net/http"
	"testing"
)

type NotRandomSource struct {
	n int
}

func (r NotRandomSource) Intn(i int) int {
	return r.n
}
func TestReturnRandomResponseSuccess(t *testing.T) {
	failRate := 30
	got := returnRandomResponse(NotRandomSource{100}, failRate)
	want := RandomResponse{http.StatusOK, SUCCESS_MESSAGE}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReturnRandomResponseFail(t *testing.T) {
	failRate := 30
	got := returnRandomResponse(NotRandomSource{1}, failRate)
	want := RandomResponse{http.StatusServiceUnavailable, FAILURE_MESSAGE}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReturnRandomResponseFailBoundary(t *testing.T) {
	failRate := 30
	got := returnRandomResponse(NotRandomSource{failRate - 1}, failRate)
	want := RandomResponse{http.StatusServiceUnavailable, FAILURE_MESSAGE}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestReturnRandomResponseSuccessBoundary(t *testing.T) {
	failRate := 30
	got := returnRandomResponse(NotRandomSource{failRate}, failRate)
	want := RandomResponse{http.StatusOK, SUCCESS_MESSAGE}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
