#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./build/main ./lambda/main.go