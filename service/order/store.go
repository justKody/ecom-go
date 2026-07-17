package order

import (
	"database/sql"

	"github.com/justKody/ecom-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) error {
	_, err := s.db.Exec("INSERT INTO orders (user_id, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}
	return nil
}
