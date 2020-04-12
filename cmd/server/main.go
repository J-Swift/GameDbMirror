package main

import (
	"os"
	"strconv"

	"github.com/J-Swift/GamesDbMirror-go/pkg/server"
)

const (
	defaultPort          = 5000
	maxResultsPerRequest = 1000
)

func resolvePort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return strconv.Itoa(defaultPort)
}

func main() {
	server.Run(resolvePort(), maxResultsPerRequest)
}
