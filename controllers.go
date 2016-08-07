package main

import (
	"github.com/kataras/iris"
	"log"
)

type OrderAPI struct {
	*iris.Context
}

type OrderItemAPI struct {
	*iris.Context
}

type TransactionAPI struct {
	 *iris.Context
}

func (request OrderAPI) Post() {
	order := Order{}
	request.ReadJSON(&order)
	//log.Print("Saving")
	order.Save()
	request.JSON(iris.StatusOK,  iris.Map{"id": order.Id})
}

func (request OrderItemAPI) Post() {
	order := Order{}

	err := order.FindId(request.Param("id"))

	if err != nil {
		log.Print(err)
		request.EmitError(iris.StatusNotFound)
		return
	}

	orderItem := OrderItem{}
	request.ReadJSON(&orderItem)

	orderItem.Save(order.Id)
	request.Text(iris.StatusOK, "")
}

func (request TransactionAPI) Post() {
	order := Order{}

	err := order.FindId(request.Param("id"))

	if err != nil {
		log.Print(err)
		request.EmitError(iris.StatusNotFound)
		return
	}

	transaction := Transaction{}
	request.ReadJSON(&transaction)

	transaction.Save(order.Id)
	request.JSON(iris.StatusOK,  iris.Map{"id": transaction.Id})
}