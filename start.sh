#!/bin/bash

echo "Starting TopV Adaptor Go..."

echo
echo "Building project..."
go mod tidy
go build -o topv-adaptor main.go push.go

echo
echo "Starting application..."
./topv-adaptor 