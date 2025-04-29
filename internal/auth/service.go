package auth

import (
	"context"
	"errors"

	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/anhilmy/tablelink-auth/pkg/grpc"
	"github.com/anhilmy/tablelink-auth/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, req *grpc.LoginRequest) (*grpc.LoginResponse, error)
	Logout(ctx context.Context) (logoutResponse, error)
	Create(ctx context.Context, req createRequest) (createResponse, error)
}

type service struct {
	repo  repository.Repository
	redis *redis.Client
}

func NewService(repo repository.Repository, rdsClient *redis.Client) Service {
	return &service{repo, rdsClient}
}

func (s *service) Login(ctx context.Context, req *grpc.LoginRequest) (*grpc.LoginResponse, error) {
	user, err := s.repo.GetUser(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := generateToken()
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	err = s.redis.Set(ctx, token, user.Name, time.Hour*24).Err()
	if err != nil {
		return nil, errors.New("failed to save token to redis")
	}

	return &grpc.LoginResponse{
		Status:  true,
		Message: "success",
		Data: &grpc.LoginResponseAccessToken{
			AccessToken: token,
		},
	}, nil
}

func (s *service) Logout(ctx context.Context) (logoutResponse, error) {
	return logoutResponse{}, nil
}

func (s *service) Create(ctx context.Context, req createRequest) (createResponse, error) {
	// TODO implement me
	panic("implement me")
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
