package payment

import (
	"context"
	"github.com/aknEvrnky/grpc-microservices-proto/golang/payment"
	"github.com/aknevrnky/grpc-microservices/order/internal/application/core/domain"
	"github.com/aknevrnky/grpc-microservices/order/internal/middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var cb = gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
	Name:        "payment_service",
	MaxRequests: 5,
	Timeout:     2 * time.Second,
	ReadyToTrip: func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	},
})

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	// Retry policy for gRPC client, it will retry 5 times if the error is Unavailable or ResourceExhausted
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
	)))

	opts = append(opts, grpc.WithUnaryInterceptor(middleware.CircuitBreakerClientInterceptor(cb)))

	// disable TLS for now
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		payment: payment.NewPaymentClient(conn),
	}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	return err
}
