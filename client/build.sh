#!/bin/bash
echo "Building RedTeam Client for multiple platforms..."
GOOS=linux GOARCH=amd64 go build -o builds/redteam-client-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o builds/redteam-client-windows-amd64.exe main.go
GOOS=darwin GOARCH=amd64 go build -o builds/redteam-client-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o builds/redteam-client-darwin-arm64 main.go
echo "Builds ready in ./builds/"