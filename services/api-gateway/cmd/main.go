package main

import (
	"log"

	"github.com/Eucastan/eucastanpay/services/gateway/internal/bootstrap"
)

func main() {

	app, err := bootstrap.New()

	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
