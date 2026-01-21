package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/dto"
	"github.com/MaiNhatHoangY2001/go-project-structure/internal/app/todoapp/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Signup(ctx context.Context, req *dto.UserSignupRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type authUsecase struct {
	userRepo      domain.UserRepository
	jwtSecret     string
	jwtExpiration int
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtSecret string, jwtExpiration int) AuthUsecase {
	return &authUsecase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (u *authUsecase) Signup(ctx context.Context, req *dto.UserSignupRequest) (*dto.UserResponse, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * time.Duration(u.jwtExpiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &dto.UserLoginResponse{
		Token: tokenString,
		User: dto.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (u *authUsecase) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(u.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
