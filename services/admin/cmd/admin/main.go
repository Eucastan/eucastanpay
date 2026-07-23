package main

import (
	"log"

	"github.com/Eucastan/eucastanpay/services/admin/internal/bootstrap"
)

func main() {

	app, err := bootstrap.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
