package main

import (
	"log"
	"net/http"

	"github.com/FrankYang0529/geekbang-golang-training-week3/server"
)

func main() {
	var srv http.Server
	app, err := server.New(&srv)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
