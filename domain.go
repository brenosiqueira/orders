package main

import (
	"github.com/satori/go.uuid"
	"time"
	"log"
	"fmt"
)

type Config struct {
	Scyllaclusters []string
	Serverport int
}

type Order struct {
	Id        string      `json:"id"`
	Number    string      `json:"number"`
	Reference string      `json:"reference"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Notes     string      `json:"notes"`
	Price     int         `json:"price"`
	Items     []OrderItem         `json:"items"`
	Transactions     []Transaction  `json:"transactions"`
}

type OrderItem struct {
	Sku       string `json:"sku"`
	UnitPrice    int    `json:"unit_price"`
	Quantity int    `json:"quantity"`
}


type Transaction struct {
	Id       string `json:"id"`
	ExternalId    string    `json:"external_id"`
	Amount int    `json:"amount"`
	Type    string `json:"type"`
	AuthorizationCode string `json:"authorization_code"`
	CardBrand string `json:"card_brand"`
	CardBin string `json:"card_bin"`
	CardLast string `json:"card_last"`
	OrderId string `json:"order_id"`
}

func (order *Order) Save() error {
	order.Id = uuid.NewV4().String()
	order.Status = "DRAFT"
	order.CreatedAt = time.Now()

	err := session.Query("INSERT INTO orders (id,number,reference,status,created_at) VALUES (?,?,?,?,?)",
		order.Id, order.Number, order.Reference, order.Status, order.CreatedAt).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}

func (order *Order) FindId(id string) error {
	return session.Query("SELECT id FROM orders WHERE id = ? ", id).Scan(&order.Id)
}

func (item *OrderItem) Save(order_id string) error {

	err := session.Query(fmt.Sprintf("UPDATE orders SET items = items + [{sku: %v, unit_price: %v, quantity: %v}] WHERE id = %v",
		item.Sku, item.UnitPrice, item.Quantity, order_id)).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}


func (order *Order) GetOrder(id string) error {

	 err := session.Query("SELECT * from orders WHERE id = ? ", id).
		 Scan(&order.Id, &order.Number, &order.Reference, &order.Status, &order.CreatedAt,
          &order.UpdatedAt, &order.Notes, &order.Price,&order.Items,&order.Transactions)

	if err != nil {
        	log.Fatal(err)
    	}

	return err
}

func (tran *Transaction) Save(order_id string) error {
	tran.Id = uuid.NewV4().String()

	query := fmt.Sprintf("UPDATE orders SET transactions = transactions + [{id: %v, external_id: '%v', amount: %v, type: '%v', authorization_code: '%v', card_brand: '%v', card_bin: '%v', card_last: '%v'}] WHERE id = %v",
		tran.Id, tran.ExternalId, tran.Amount, tran.Type, tran.AuthorizationCode, tran.CardBrand, tran.CardBin, tran.CardLast, order_id)

	err := session.Query(query).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}
