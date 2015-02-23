#!/bin/bash
c=1
while [ $c -le 5 ]
do 
	`go run text.go`
	sleep 1
done