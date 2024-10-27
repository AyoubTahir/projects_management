package services

import (
	"context"

	"github.com/AyoubTahir/projects_management/internal/repositories"
	"github.com/AyoubTahir/projects_management/pkg/types"
)

type Service struct {
	repository *repositories.Repository
	User       UserServiceI
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{
		repository: repository,
		User:       NewUserService(repository),
	}
}

type UserServiceI interface {
	CreateUser(ctx context.Context, user *types.CreateUserPayload) (map[string]interface{}, error)
	GetUserByID(ctx context.Context, id int64) (map[string]interface{}, error)
	// Add other user-related methods as needed
}
