package main

import (
	"flag"
	"log"

	"funny/fastapi"
	"funny/fastapi/example/module1"

	"github.com/funny/fastbin"
)

func main() {
	gencode := flag.Bool("gencode", false, "generate code")
	flag.Parse()

	app := fastapi.New()
	app.Register(1, &module1.Service{})

	if *gencode {
		fastapi.GenCode(app)
		fastbin.GenCode()
		return
	}

	server, err := app.Listen("tcp", "0.0.0.0:0", nil)
	if err != nil {
		log.Fatal("setup server failed:", err)
	}
	go server.Serve()

	client.
}
