package controller

import (
	"L0/service"
	"encoding/json"
	"log"
	"net/http"
)

type OrderController struct {
	OrderService *service.OrderService
}

func NewOrderController(service *service.OrderService) *OrderController {
	return &OrderController{
		OrderService: service,
	}
}

func (c *OrderController) MainPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal("main", err.Error())
	}
	id := r.Form.Get("id")
	if id == "" {
		w.Write([]byte("hello\n"))
		return
	}

	order, ok := c.OrderService.GetById(id)
	if ok {
		response, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err.Error())
		}

		w.Write(response)
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Такого заказа не существует!"))
	}
}
