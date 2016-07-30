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
