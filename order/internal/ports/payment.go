package ports

import "github.com/aknevrnky/microservices-order/internal/application/core/domain"

type PaymentPort interface {
	Charge(order *domain.Order) error
}
