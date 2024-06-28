package middleware

import (
	"context"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
)

func CircuitBreakerClientInterceptor(cb *gobreaker.CircuitBreaker[any]) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (any, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)

			if err != nil {
				return nil, err
			}

			return reply, nil
		})

		return cbErr
	}
}
