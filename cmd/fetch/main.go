package main

import "github.com/J-Swift/GamesDbMirror-go/pkg/fetch"

const (
	dataDir = ".fetch-cache"
)

func main() {
	fetch.Run(dataDir)
}
