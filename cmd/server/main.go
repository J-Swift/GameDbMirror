package main

import (
	"os"
	"strconv"

	"github.com/J-Swift/GamesDbMirror-go/pkg/server"
)

const (
	defaultPort          = 5000
	maxResultsPerRequest = 1000
	dataDir              = ".fetch-cache"
)

func resolvePort() string {
	// NOTE(jpr): heroku requires you to bind to the port they specify through the envvar
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return strconv.Itoa(defaultPort)
}

func main() {
	server.Run(dataDir, resolvePort(), maxResultsPerRequest)
}
