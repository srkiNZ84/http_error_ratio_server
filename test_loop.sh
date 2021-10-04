#!/bin/sh

while :
do
  curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/
  echo
  sleep 1
done
