package main

import (
	"github.com/satori/go.uuid"
	"time"
	"log"
)

type Order struct {
	Id        string      `json:"id"`
	Number    string      `json:"number"`
	Reference string      `json:"reference"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Notes     string      `json:"notes"`
	Price     int         `json:"price"`
}

type OrderItem struct {
	Sku       string `json:"sku"`
	UnitPrice    int    `json:"unit_price"`
	Quantity int    `json:"quantity"`
	OrderId    string `json:"order_id"`
}


func (order *Order) Save() error {
	//log.Print("Saving to disk")
	order.Id = uuid.NewV4().String()
	order.Status = "DRAFT"
	order.CreatedAt = time.Now()

	err := session.Query("INSERT INTO \"order\" (id,number,reference,status,created_at) VALUES (?,?,?,?,?)",
		order.Id, order.Number, order.Reference, order.Status, order.CreatedAt).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}

func (order *Order) FindId(id string) error {
	return session.Query("SELECT id FROM \"order\" WHERE id = ? ", id).Scan(&order.Id)
}

func (item *OrderItem) Save(order_id string) error {
	err := session.Query("INSERT INTO \"order_item\" (sku,order_id,unit_price,quantity) VALUES (?,?,?,?)",
		item.Sku, order_id, item.UnitPrice, item.Quantity).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}