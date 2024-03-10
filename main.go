package main

import (
	"L0/entity"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer nc.Close()

	var data []byte
	order := entity.Order{}

	nc.Subscribe("order", func(m *nats.Msg) {
		data = m.Data
		err = json.Unmarshal(data, &order)
		if err != nil {
			log.Fatal(err.Error())
		}

		conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/L0")
		if err != nil {
			log.Fatal(err.Error())
		}

		var jsonOrder []byte
		jsonOrder, err = json.Marshal(order)
		if err != nil {
			log.Fatal(err.Error())
		}

		_, err = conn.Exec(context.Background(), "insert into orders (id, data) values ($1, $2)", order.OrderUid, jsonOrder)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
