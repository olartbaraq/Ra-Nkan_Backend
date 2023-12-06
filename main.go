package main

import (
	"github.com/olartbaraq/spectrumshelf/api"
)

func main() {
	server := api.NewServer(".")
	server.Start(8000)
}
