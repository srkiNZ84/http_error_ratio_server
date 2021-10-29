#!/bin/sh

hey -n 100000 http://localhost:8080/

# -t 0 infinite timeout
# -c 2 two concurrent clients
# -n total number of requests to make
hey -t 0 -c 2 -n 20 http://localhost:8080/