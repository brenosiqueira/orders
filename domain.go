package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"time"
	"errors"
	"github.com/gocql/gocql"
	"strings"
)

type Config struct {
	Scyllaclusters []string
	Serverport     int
}

type Order struct {
	Id           string        `json:"id"`
	Number       string        `json:"number"`
	Reference    string        `json:"reference"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Notes        string        `json:"notes"`
	Price        int           `json:"price"`
	Items        []OrderItem   `json:"items"`
	Transactions []Transaction `json:"transactions"`
}

func (o Order) ValidadeNewOrder() error {
	if o.Number == "" {
		return errors.New("Order number cannot be empty")
	}

	if o.Reference == "" {
		return errors.New("Order Reference cannot be empty")
	}

	if o.Status == "" {
		return errors.New("Order Status cannot be empty")
	}

	if o.Price <= 0 {
		return errors.New("Order Price cannot be less or equal to 0")
	}

	if len(o.Items) > 0 {
		for i := range o.Items {
			return o.Items[i].ValidadeNewOrderItem();
		}
	}

	if len(o.Transactions) > 0 {
		for k := range o.Transactions {
			return o.Transactions[k].ValidateNewTransaction()
		}
	}
	return nil
}

type OrderItem struct {
	Sku       string `json:"sku"`
	UnitPrice int    `json:"unit_price"`
	Quantity  int    `json:"quantity"`
}


func (oi OrderItem) ValidadeNewOrderItem() error {
	if oi.Quantity <= 0 {
		return errors.New("OrderItem Quantity cannot be less or equal to 0")
	}

	if oi.Sku == "" {
		return errors.New("OrderItem Sku cannot be empty")
	}

	if oi.UnitPrice < 0 {
		return  errors.New("OrderItem UnitPrice cannot be less than 0")
	}

	return nil;
}


type Transaction struct {
	Id                string `json:"id"`
	ExternalId        string `json:"external_id"`
	Amount            int    `json:"amount"`
	Type              string `json:"type"`
	AuthorizationCode string `json:"authorization_code"`
	CardBrand         string `json:"card_brand"`
	CardBin           string `json:"card_bin"`
	CardLast          string `json:"card_last"`
	OrderId           string `json:"order_id"`
}

func (t Transaction) ValidateNewTransaction() error {

	if t.ExternalId == "" {
		return errors.New("Transaction ExternalId cannot be empty")
	}

	if t.Amount <= 0 {
		return  errors.New("Transaction Amount cannot be less or equal to 0")
	}

	if t.AuthorizationCode == "" {
		return errors.New("Transaction Authorization code cannot be empty")
	}

	if t.Type == "" {
		return errors.New("Transaction Type cannot be empty")
	}

	if t.CardBin == "" {
		return  errors.New("Transaction Card Bin cannot be empty")
	}

	if t.CardBrand == "" {
		return errors.New("Transaction Card Brand cannot be empty")
	}

	if t.CardLast == "" {
		return errors.New("Transaction Card Last cannot be empty")
	}

	if t.OrderId == "" {
		return errors.New("Transaction OrderId cannot be empty")
	}

	return nil;
}


func (order *Order) Save() error {
	order.Id = uuid.NewV4().String()
	order.Status = "DRAFT"
	order.CreatedAt = time.Now()

	err := order.ValidadeNewOrder()
	if err != nil {
		log.Print(err)
		return err;
	}

	err = session.Query("INSERT INTO orders (id,number,reference,status,created_at) VALUES (?,?,?,?,?)",
		order.Id, order.Number, order.Reference, order.Status, order.CreatedAt).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

func (order *Order) FindId(id string) error {
	return session.Query("SELECT id FROM orders WHERE id = ? ", id).Scan(&order.Id)
}

func (item *OrderItem) Save(order_id string) error {


	err := item.ValidadeNewOrderItem()
	if err != nil {
		log.Print(err)
		return err;
	}

	err = session.Query(fmt.Sprintf("UPDATE orders SET items = items + [{sku: %v, unit_price: %v, quantity: %v}] WHERE id = %v",
		item.Sku, item.UnitPrice, item.Quantity, order_id)).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

func (order *Order) GetOrder(id string) error {
	err := session.Query("SELECT id, number, reference, status, notes, price from orders WHERE id = ? ", id).
		 Scan(&order.Id, &order.Number, &order.Reference, &order.Status, &order.Notes, &order.Price)

	if err != nil {
		log.Print(err)
	}

	return err
}

func (tran *Transaction) Save(order_id string) error {
	tran.Id = uuid.NewV4().String()

	err := tran.ValidateNewTransaction()
	if err != nil {
		log.Print(err)
		return err;
	}
	query := fmt.Sprintf("UPDATE orders SET transactions = transactions + [{id: %v, external_id: '%v', amount: %v, type: '%v', authorization_code: '%v', card_brand: '%v', card_bin: '%v', card_last: '%v'}] WHERE id = %v",
		tran.Id, tran.ExternalId, tran.Amount, tran.Type, tran.AuthorizationCode, tran.CardBrand, tran.CardBin, tran.CardLast, order_id)

	err = session.Query(query).Exec()

	if err != nil {
		log.Print(err)
	}

	return err
}

func (order *OrderItem) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {

	t := strings.SplitN(string(data), " ", 2)

	log.Print(t)
	return nil

}
