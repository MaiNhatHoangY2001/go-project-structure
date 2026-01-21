package grpc

import (
	"context"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/usecase"
	pb "github.com/MaiNhatHoangY2001/go-project-structure/proto/auth"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authUsecase usecase.AuthUsecase
}

func NewAuthServer(authUsecase usecase.AuthUsecase) *AuthServer {
	return &AuthServer{
		authUsecase: authUsecase,
	}
}

func (s *AuthServer) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	userReq := &dto.UserSignupRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	user, err := s.authUsecase.Signup(ctx, userReq)
	if err != nil {
		return nil, err
	}

	return &pb.SignupResponse{
		UserId:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginReq := &dto.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := s.authUsecase.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token:  result.Token,
		UserId: result.User.ID,
		Email:  result.User.Email,
		Name:   result.User.Name,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := s.authUsecase.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: err.Error(),
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: "invalid token claims",
		}, nil
	}

	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
	}, nil
}
