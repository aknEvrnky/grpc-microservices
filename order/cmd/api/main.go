package main

import (
	"github.com/aknevrnky/microservices-order/config"
	"github.com/aknevrnky/microservices-order/internal/adapters/db"
	"github.com/aknevrnky/microservices-order/internal/adapters/grpc"
	"github.com/aknevrnky/microservices-order/internal/adapters/payment"
	"github.com/aknevrnky/microservices-order/internal/application/core/api"
	"log"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Error while creating db adapter: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Error while creating payment adapter: %v", err)
	}

	app := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(app, config.GetApplicationPort())
	grpcAdapter.Run()
}
