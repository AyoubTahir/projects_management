package repositories

import (
	"context"

	"github.com/AyoubTahir/projects_management/pkg/orm"
	"github.com/AyoubTahir/projects_management/pkg/types"
)

type Repository struct {
	orm  *orm.Orm
	User UserRepositoryI
}

func NewRepository(orm *orm.Orm) *Repository {
	return &Repository{
		orm:  orm,
		User: NewUserRepository(orm),
		// Initialize OrderRepository here when you have it
	}
}

type UserRepositoryI interface {
	Create(ctx context.Context, user *types.CreateUserPayload) (map[string]interface{}, error)
	GetByID(ctx context.Context, id int64) (map[string]interface{}, error)
	// Add other user-related methods as needed
}
