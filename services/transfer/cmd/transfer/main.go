// Package main Transfer Service API
//
// @title           EucastanPay Transfer Service API
// @version         1.0
// @description     Authentication and Transfer Management Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host transfer-sby1.onrender.com
// @BasePath /api/v1
// @schemes https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer <JWT>
package main

import (
	"log"

	"github.com/Eucastan/eucastanpay/services/transfer/internal/bootstrap"
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
