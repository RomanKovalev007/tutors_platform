package auth

import (
	"api_gateway/pkg/resilience"
	"context"
	"fmt"
	"time"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn           *grpc.ClientConn
	service        pb.AuthServiceClient
	circuitBreaker *resilience.CircuitBreaker
	retry          *resilience.Retry
}

func NewClient(address string, opts ...grpc.DialOption) (*Client, error) {
	cbConfig := resilience.CircuitBreakerConfig{
		Name:             "auth-service",
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
		OnStateChange:    resilience.LoggingStateChangeCallback,
	}

	retryConfig := resilience.RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     2 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}

	resilienceInterceptor := resilience.NewResilienceInterceptor(cbConfig, retryConfig, 10*time.Second)

	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(resilienceInterceptor.UnaryClientInterceptor()),
	}

	if len(opts) == 0 {
		opts = defaultOpts
	} else {
		opts = append(opts, defaultOpts...)
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return &Client{
		conn:           conn,
		service:        pb.NewAuthServiceClient(conn),
		circuitBreaker: resilience.NewCircuitBreaker(cbConfig),
		retry:          resilience.NewRetry(retryConfig),
	}, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (string, error) {
	resp, err := c.service.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return "", err
	}
	if !resp.IsValid {
		return "", fmt.Errorf("token invalid")
	}
	return resp.UserId, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
