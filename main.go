package main

import (
	"L0/controller"
	"L0/entity"
	"L0/service"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
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

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/L0")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close(context.Background())

	orderService := service.NewOrderService(conn)
	orderController := controller.NewOrderController(orderService)
	orderService.FillCache()

	subscribe(nc, conn, orderService)

	http.HandleFunc("/", orderController.MainPage)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func subscribe(nc *nats.Conn, conn *pgx.Conn, orderService *service.OrderService) {
	nc.Subscribe("order", func(m *nats.Msg) {
		data := m.Data                  // принятое сообщение
		acceptedOrder := entity.Order{} // принятый объект типа Order
		err := json.Unmarshal(data, &acceptedOrder)

		if err != nil {
			log.Println(err.Error())
		}

		_, err = conn.Exec(context.Background(),
			"insert into orders (id, data) values ($1, $2)", acceptedOrder.OrderUid, acceptedOrder)

		if err != nil {
			log.Println(err.Error())
		}
		orderService.AddOrder(acceptedOrder.OrderUid, acceptedOrder) // добавление в словарь принятого заказа
	})
}
