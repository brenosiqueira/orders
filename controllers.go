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

func (request OrderAPI) Post() {
	order := Order{}
	request.ReadJSON(&order)
	//log.Print("Saving")
	order.Save()
	request.JSON(iris.StatusOK,  iris.Map{"id": order.Id})
	//request.JSON(iris.StatusOK, order.Id)
}

func (request OrderItemAPI) Post() {
	order := Order{}

	err := order.Find(request.Param("id"))
	if err != nil {
		log.Print(err)
		request.EmitError(iris.StatusNotFound)
		return
	}

	item := OrderItem{}
	request.ReadJSON(&item)

	log.Print(request.Params)


	log.Print(order)

	item.Save(order.Id)
	request.Text(iris.StatusOK, "")
	//order := Order.Find(request.Param("id"))
	//order := Order{}.Find(request.Param("id"))
	//order.find(request.Param("id"))


}