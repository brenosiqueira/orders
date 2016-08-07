package main

import (
	"github.com/satori/go.uuid"
	"time"
	"log"
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


func (orderDetails *OrderDetails) GetOrder(id string) error {

	 err := session.Query("SELECT \"order\".id,\"order\".number, \"order\".reference, \"order\".status, \"order\".created_at," +
 	  											"\"order\".updated_at, \"order\".notes, \"order_item\".sku, \"order_item\".unit_price, \"order_item\".quantity," +
													"\"transaction\".external_id, \"transaction\".type, \"transaction\".amount,"+
 													"\"transaction\".authorization_code, \"transaction\".card_brand, \"transaction\".card_bin, \"transaction\".card_last"+
 													" FROM \"order\" "+
 													" RIGHT JOIN \"order_item\" on \"order\".id = \"order_item\".order_id" +
 													" RIGHT JOIN \"transaction\" \"order\".id = \"transaction\".order_id WHERE id = ? ", id).
		 Scan(&orderDetails.Id, &orderDetails.Number, &orderDetails.Reference, &orderDetails.Status, &orderDetails.CreatedAt,
											&orderDetails.UpdatedAt, &orderDetails.Notes, &orderDetails.Sku, &orderDetails.UnitPrice, &orderDetails.Quantity,
										  &orderDetails.ExternalId, &orderDetails.Type, &orderDetails.Amount, &orderDetails.AuthorizationCode, &orderDetails.CardBrand,
	  									&orderDetails.CardBin, &orderDetails.CardLast)

	if err != nil {
        	log.Fatal(err)
    	}

	return err
}

func (tran *Transaction) Save(order_id string) error {
	tran.Id = uuid.NewV4().String()

	err := session.Query("INSERT INTO \"transaction\" (id,order_id,external_id,amount,type,authorization_code,card_brand,card_bin,card_last) VALUES (?,?,?,?,?,?,?,?,?)",
		tran.Id, order_id, tran.ExternalId, tran.Amount, tran.Type, tran.AuthorizationCode, tran.CardBrand, tran.CardBin, tran.CardLast).Exec()

	if (err != nil) {
		log.Fatal(err)
	}

	return err
}
