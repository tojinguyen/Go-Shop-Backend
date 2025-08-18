package adapter

import (
	"context"
	"log"

	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	user_v1 "github.com/toji-dev/go-shop/proto/gen/go/user/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceAdapter interface {
	GetAddressById(ctx context.Context, addressID string) (*user_v1.Address, error)
	Close() error
}

type grpcUserAdapter struct {
	conn   *grpc.ClientConn
	client user_v1.UserServiceClient
}

func NewGrpcUserAdapter(userServiceAddr string) (UserServiceAdapter, error) {
	log.Printf("Connecting to user service at %s", userServiceAddr)
	conn, err := grpc.NewClient(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
		return nil, err
	}

	client := user_v1.NewUserServiceClient(conn)

	log.Printf("Successfully connected to user service at %s", userServiceAddr)

	return &grpcUserAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcUserAdapter) GetAddressById(ctx context.Context, addressID string) (*user_v1.Address, error) {
	log.Printf("Requesting address with ID: %s", addressID)
	req := &user_v1.GetAddressRequest{
		AddressId: addressID,
	}

	resp, err := a.client.GetAddressById(ctx, req)

	if err != nil {
		log.Printf("Error fetching address by ID: %v", err)
		return nil, err
	}
	if resp == nil || resp.Address == nil {
		log.Printf("No address found for ID: %s", addressID)
		return nil, apperror.NewNotFound("Address", addressID)
	}

	log.Printf("Successfully fetched address with ID: %v", resp.Address)

	return resp.Address, nil
}

func (a *grpcUserAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
