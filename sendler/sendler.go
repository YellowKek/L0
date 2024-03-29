package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("model.json")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer nc.Close()

	nc.Publish("order", data)
}
