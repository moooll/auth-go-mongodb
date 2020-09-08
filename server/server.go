package main

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/gobuffalo/uuid"
	"github.com/valyala/fasthttp"
)

var ctx *fasthttp.RequestCtx

//IDs slice stores uuids of all users created for this app. (3 in this case)
var IDs []uuid.UUID

func seedUsers() {
	for i := 0; i < 3; i++ {
		IDs = append(IDs, uuid.Must(uuid.NewV4()))
		id := writeUser(IDs[i])
		if id == nil {
			log.Fatal("no doc in the DB")
		}
	}
}

func main() {
	router := router.New()
	connectToTheDB()
	seedUsers()
	router.GET("/", home)
	router.GET("/setTokens", setTokens)
	router.POST("/refresh", refresh)
	router.POST("/delOne", delOne)
	router.POST("/delAll", delAll)

	log.Fatal(fasthttp.ListenAndServe(":8084", router.Handler))
}
