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

	err := order.Save()

	if err != nil {
		log.Print(err)
		request.Text(iris.StatusBadRequest, err.Error())
		return
	}

	request.JSON(iris.StatusOK, iris.Map{"id": order.Id})
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

	err = orderItem.Save(order.Id)

	if err != nil {
		log.Print(err)
		request.Text(iris.StatusBadRequest, err.Error())
		return
	}
	request.Text(iris.StatusOK, "")
}

func (request OrderAPI) Get() {
	order := Order{}

	err := order.GetOrder(request.Param("id"))


	if err != nil {
		log.Print(err)
		request.EmitError(iris.StatusNotFound)
		return
	}

	request.ReadJSON(&order)
	request.JSON(iris.StatusOK, order)
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

	err = transaction.Save(order.Id)

	if err != nil {
		request.Text(iris.StatusBadRequest,err.Error())
	}
	request.JSON(iris.StatusOK, iris.Map{"id": transaction.Id})
}
