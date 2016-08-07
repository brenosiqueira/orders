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

type OrderDetails struct {
	Id        string      				`json:"id"`
	Number    string      				`json:"number"`
	Reference string      				`json:"reference"`
	Status    string      				`json:"status"`
	CreatedAt time.Time   				`json:"createdAt"`
	UpdatedAt time.Time   				`json:"updatedAt"`
	Notes     string      				`json:"notes"`
	Price     int         				`json:"price"`
	Sku       string 							`json:"sku"`
	UnitPrice    int    					`json:"unit_price"`
	Quantity int    							`json:"quantity"`
	ExternalId string    				`json:"externalId"`
	Amount int      							`json:"amount"`
	Type string     							`json:"type"`
	AuthorizationCode string 		`json:"authorizationCode"`
	CardBrand string  						`json:"cardBrand"`
  CardBin string 							`json:"cardBin"`
	CardLast string 							`json:"cardLast"`
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

func (order *Order) Find(id string) error {
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

func (orderDetails *OrderDetails) GetOrder(id string) error {

	 err := session.QueryRow("SELECT \"order\".id,\"order\".number, \"order\".reference, \"order\".status, \"order\".created_at,
 	  											\"order\".updated_at, \"order\".notes, \"order_item\".sku, \"order_item\".unit_price, \"order_item\".quantity,
													\"transaction\".external_id, \"transaction\".type, \"transaction\".amount,
 													\"transaction\".authorization_code, \"transaction\".card_brand, \"transaction\".card_bin, \"transaction\".card_last
 													FROM \"order\"
 													RIGHT JOIN \"order_item\" on \"order\".id = \"order_item\".order_id
 													RIGHT JOIN \"transaction\" \"order\".id = \"transaction\".order_id WHERE id = ? ", id)
	  						.Scan(&orderDetails.Id, &orderDetails.number, &orderDetails.Reference, &orderDetails.Status, &orderDetails.CreatedAt,
											&orderDetails.UpdatedAt, &orderDetails.Notes, &orderDetails.Sku, &orderDetails.UnitPrice, &orderDetails.Quantity,
										  &orderDetails.ExternalId, &orderDetails.Type, &orderDetails.Amount, &orderDetails.AuthorizationCode, &orderDetails.CardBrand,
	  									&orderDetails.CardBin, &orderDetails.CardLast)

		if err != nil {
        return &OrderDetails{}, err
    } else {
        return &orderDetails, nil
    }

}
