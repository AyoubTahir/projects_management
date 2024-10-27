package services

import (
	"context"
	"fmt"

	"github.com/AyoubTahir/projects_management/internal/repositories"
	"github.com/AyoubTahir/projects_management/pkg/types"
)

type UserService struct {
	repository *repositories.Repository
}

func NewUserService(repository *repositories.Repository) UserServiceI {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(ctx context.Context, user *types.CreateUserPayload) (map[string]interface{}, error) {
	data, err := s.repository.User.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return data, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (map[string]interface{}, error) {
	user, err := s.repository.User.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}
