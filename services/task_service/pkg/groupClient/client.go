package client

import (
	"context"
	"fmt"
	"time"

	"task_service/pkg/resilience"

	pb "github.com/RomanKovalev007/tutors_platform/api/gen/go/group"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GroupClient struct {
	conn           *grpc.ClientConn
	client         pb.GroupsServiceClient
	timeout        time.Duration
	circuitBreaker *resilience.CircuitBreaker
	retry          *resilience.Retry
}

func NewGroupClient(address string) (*GroupClient, error) {
	cbConfig := resilience.CircuitBreakerConfig{
		Name:             "group-service",
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

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(resilienceInterceptor.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	client := pb.NewGroupsServiceClient(conn)

	conn.Connect()

	return &GroupClient{
		conn:           conn,
		client:         client,
		timeout:        5 * time.Second,
		circuitBreaker: resilience.NewCircuitBreaker(cbConfig),
		retry:          resilience.NewRetry(retryConfig),
	}, nil
}

func (c *GroupClient) GetGroupInfo(ctx context.Context, groupID string) (*pb.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.GetGroupRequest{
		Id:             groupID,
		IncludeMembers: false,
	}

	resp, err := c.client.GetGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetGroup(), nil
}

func (c *GroupClient) GetGroupMembers(ctx context.Context, groupID string) ([]*pb.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.ListGroupMembersRequest{
		GroupId: groupID,
	}

	resp, err := c.client.ListGroupMembers(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetMembers(), nil
}

func (c *GroupClient) Close() error {
	return c.conn.Close()
}
