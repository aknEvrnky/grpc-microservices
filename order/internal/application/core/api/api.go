package api

import (
	"github.com/aknevrnky/grpc-microservices/order/internal/application/core/domain"
	"github.com/aknevrnky/grpc-microservices/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a *Application) PlaceOrder(order *domain.Order) (*domain.Order, error) {
	err := a.db.Save(order)
	if err != nil {
		return &domain.Order{}, err
	}

	paymentErr := a.payment.Charge(order)
	// if payment failed, we need to return the error gracefully
	if paymentErr != nil {
		// Converts a complex error to a status
		st := status.Convert(paymentErr)
		// Slices for whole errors
		var allErrors []string
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					allErrors = append(allErrors, violation.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrors, "\n"),
		}

		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)

		return &domain.Order{}, statusWithDetails.Err()
	}

	return order, nil
}
