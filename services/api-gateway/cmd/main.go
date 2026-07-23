// @title           EucastanPay API Gateway
// @version         1.0
// @description     API Gateway for the EucastanPay Microservices Platform.
//
// @contact.name    Stanley Emeh
//
// @BasePath        /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	_ "github.com/Eucastan/eucastanpay/services/api-gateway/docs"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/bootstrap"
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
