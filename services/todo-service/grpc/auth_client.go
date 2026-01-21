package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(host string, port string) (*AuthClient, error) {
	address := fmt.Sprintf("%s:%s", host, port)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client: pb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (string, error) {
	req := &pb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return "", err
	}

	if !resp.Valid {
		return "", fmt.Errorf("invalid token: %s", resp.ErrorMessage)
	}

	return resp.UserId, nil
}

func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
