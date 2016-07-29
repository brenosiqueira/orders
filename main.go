package main

import (
	"github.com/gocql/gocql"
	"github.com/kataras/iris"
	_ "github.com/dimiro1/banner/autoload"
	"log"
)

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "orders"
	cluster.Consistency = gocql.One

	session, _ := cluster.CreateSession()
	defer session.Close()

	setupWebServer(session);
}

func setupWebServer(session *gocql.Session) {
	iris.Get("/hello", hello)
	iris.Get("/scylla", func(ctx *iris.Context) {
		scylla(ctx, session)
	})

	iris.Listen(":8080")
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

