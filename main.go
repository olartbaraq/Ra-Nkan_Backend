package main

import (
	"github.com/olartbaraq/spectrumshelf/api"
)

func main() {

	//rabbitmq.RabbitMQServer()
	server := api.NewServer(".")
	server.Start(8000)
}
