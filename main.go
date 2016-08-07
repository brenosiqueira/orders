package main

import (
	"github.com/gocql/gocql"
	"github.com/kataras/iris"
	_ "github.com/dimiro1/banner/autoload"
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
	"fmt"
)

var session *gocql.Session
var config Config

func main() {
	config = Config{}
	configFile, err := ioutil.ReadFile("application.yml")
	if err != nil {
		panic(err)
	}

	log.Print(config)

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	log.Print(config)

	cluster := gocql.NewCluster(config.Scyllaclusters...)
	cluster.Keyspace = "orders"
	cluster.Consistency = gocql.One

	session, _ = cluster.CreateSession()
	defer session.Close()

	setupWebServer(session);
}

func setupWebServer(session *gocql.Session) {
	iris.Get("/hello", hello)
	iris.Get("/scylla", func(ctx *iris.Context) {
		scylla(ctx, session)
	})

	iris.API("/orders", OrderAPI{})
	iris.API("/orders/:id", OrderDetailsAPI{})
	iris.API("/orders/:id/items", OrderItemAPI{})
	iris.API("/orders/:id/transactions", TransactionAPI{})

	iris.Listen(fmt.Sprintf(":%v",config.Serverport))
}

func scylla(ctx *iris.Context, session *gocql.Session) {
	if err := session.Query("INSERT INTO \"order\" (id, number) VALUES (?, ?)",
		gocql.TimeUUID(), "10").Exec(); err != nil {
		log.Fatal(err)
	}

	ctx.Write("OK")
}

func hello(ctx *iris.Context) {
	ctx.Write("World!\n")
}

