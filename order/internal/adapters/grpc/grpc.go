package grpc

import (
	"context"
	"github.com/aknEvrnky/grpc-microservices-proto/golang/order"
	"github.com/aknevrnky/microservices-order/internal/application/core/domain"
	"log"
)

func (a *Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range request.Items {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.Name,
			UnitPrice:   float64(orderItem.UnitPrice),
			Quantity:    int32(orderItem.Quantity),
		})
	}

	newOrder := domain.NewOrder(request.UserId, orderItems)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		log.Fatalf("failed to place order, error: %v\n", err)
	}

	return &order.CreateOrderResponse{
		OrderId: result.ID,
	}, nil
}
