package service

import (
	"context"
	"time"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/models"
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	AddUser(ctx context.Context, user models.User) error
	AddRefreshToken(ctx context.Context, email string, refreshTokenHash string, ttl time.Duration, ipAddress string) error
	GetRefreshTokenProps(ctx context.Context, email string) (models.RefreshToken, error)
}

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

type EmailRepository interface {
	NotifyUser(ctx context.Context, email string) error
}

type EmailService struct {
	Repo EmailRepository
}

func NewEmailService(repo EmailRepository) *EmailService {
	return &EmailService{
		Repo: repo,
	}
}
