# HTTP error ratio server
[![Go Report Card](https://goreportcard.com/badge/github.com/srkiNZ84/http_error_ratio_server)](https://goreportcard.com/report/github.com/srkiNZ84/http_error_ratio_server)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/srkiNZ84/http_error_ratio_server)
![unit tests](https://github.com/srkiNZ84/http_error_ratio_server/actions/workflows/test.yml/badge.svg)


## About

A trivial web server that returns (on average) a set ratio of errors (defaults to 30% of responses being
errors) on requests. The main use case is testing monitors and alerts that rely on success/fail ratios.

## How to build

To build the web server, make sure you have the `go` command installed and then run:

```
go build main.go
```

## How to run

To run the web server locally, make sure you have the `go` command installed and then run:

```
go run main.go
```

Alternatively, just run the binary created as a part of the build step.

## Configuration

The applicaiton has two configuration variables, both of which are read in from the environment variables of the same name: 

* PORT - The port that the HTTP server should listen on (defaults to TCP port 8080)

* ERROR_RATIO - The approximate percentage of requests (represented by an integer between 1 and 100) that should return an HTTP error (defaults to 30)

## Docker

### Run Docker image

```
docker run -p 8080:8080 -e ERROR_RATIO=30 -e PORT=8090 srkinz/http-error-ratio
```

### Build

To build the docker image, run:

```
docker build .
```