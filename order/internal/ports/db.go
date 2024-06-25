package ports

import "github.com/aknevrnky/microservices-order/internal/application/core/domain"

type DBPort interface {
	Get(id int64) (domain.Order, error)
	Save(order *domain.Order) error
}
