package service

import (
	"L0/entity"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"log"
)

type OrderService struct {
	conn   *pgx.Conn
	orders map[string]entity.Order
}

func NewOrderService(conn *pgx.Conn) *OrderService {
	return &OrderService{
		conn:   conn,
		orders: map[string]entity.Order{},
	}
}

func (s *OrderService) FillCache() {
	rows, err := s.conn.Query(context.Background(), "select data from orders")
	if err != nil {
		log.Fatal(err.Error())
	}
	for rows.Next() {
		var d []byte
		if err := rows.Scan(&d); err != nil {
			log.Fatal(err.Error())
		}
		order := entity.Order{}
		if err := json.Unmarshal(d, &order); err != nil {
			log.Fatal(err.Error())
		}
		s.AddOrder(order.OrderUid, order)
	}
}

func (s *OrderService) GetById(id string) (entity.Order, bool) {
	order, ok := s.orders[id]
	return order, ok
}

func (s *OrderService) AddOrder(id string, order entity.Order) {
	s.orders[id] = order
}
