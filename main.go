package main

import (
	"github.com/locpham24/go-shorten-url/db"
	"github.com/locpham24/go-shorten-url/handler"
)
func main() {
	redisDb := db.RedisDb{}
	redisDb.Connect()

	router := handler.InitRouter(redisDb.Client)
	router.Run(":8080")
}


