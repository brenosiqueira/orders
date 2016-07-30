package main

import (
	"github.com/kataras/iris"
	//"log"
)

type OrderAPI struct {
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
