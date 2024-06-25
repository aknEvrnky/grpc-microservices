package ports

import "github.com/aknevrnky/grpc-microservices/order/internal/application/core/domain"

type ApiPort interface {
	PlaceOrder(order *domain.Order) (*domain.Order, error)
}
