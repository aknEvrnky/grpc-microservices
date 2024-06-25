package payment

import (
	"context"
	"github.com/aknEvrnky/grpc-microservices-proto/golang/payment"
	"github.com/aknevrnky/grpc-microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &Adapter{
		client: payment.NewPaymentClient(conn),
	}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	req := payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	}

	_, err := a.client.Create(context.Background(), &req)

	return err
}
