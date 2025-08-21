package client

import (
	"context"

	"github.com/uchupx/saceri-chatbot-api/pkg/grpc/proto/gen/authservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client authservice.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AuthClient{conn: conn, client: authservice.NewAuthServiceClient(conn)}, nil
}

func (ac *AuthClient) GetUser(ctx context.Context, token string) (*authservice.GetUserResponse, error) {
	payload := authservice.GetUserRequest{
		Token: token,
	}

	res, err := ac.client.GetUser(ctx, &payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ac *AuthClient) Register(ctx context.Context, payload *authservice.RegisterUserRequest) (*authservice.RegisterUserResponse, error) {
	res, err := ac.client.RegisterUser(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ac *AuthClient) Login(ctx context.Context, payload *authservice.LoginRequest) (*authservice.LoginResponse, error) {
	res, err := ac.client.Login(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}
