package adapter

import (
	"context"

	user_v1 "github.com/toji-dev/go-shop/proto/gen/go/user/v1"
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
	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := user_v1.NewUserServiceClient(conn)

	return &grpcUserAdapter{
		conn:   conn,
		client: client,
	}, nil
}

func (a *grpcUserAdapter) GetAddressById(ctx context.Context, addressID string) (*user_v1.Address, error) {
	req := &user_v1.GetAddressRequest{
		AddressId: addressID,
	}

	resp, err := a.client.GetAddressById(ctx, req)

	if err != nil {
		return nil, err
	}
	return resp.Address, nil
}

func (a *grpcUserAdapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
