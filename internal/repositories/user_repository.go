package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AyoubTahir/projects_management/pkg/orm"
	"github.com/AyoubTahir/projects_management/pkg/types"
)

type UserRepository struct {
	orm *orm.Orm
}

func NewUserRepository(orm *orm.Orm) UserRepositoryI {
	return &UserRepository{orm: orm}
}

func (r *UserRepository) Create(ctx context.Context, user *types.CreateUserPayload) (map[string]interface{}, error) {
	//query := `INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	data, err := r.orm.Table("users").Create(map[string]interface{}{
		"username": user.UserName,
		"email":    user.Email,
		"password": user.Password,
	})
	//err := r.db.QueryRowContext(ctx, query, user.UserName, user.Email, user.Password, time.Now(), time.Now()).Scan(&user.UserName)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return data, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (map[string]interface{}, error) {
	//query := `SELECT id, name, email, password FROM users WHERE id = $1`
	//user := &models.User{}
	//err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	data, err := r.orm.Table("users").
		WithContext(ctx).
		Select("id", "username", "email", "password").
		Where("id", "=", id).
		First()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return data, nil
}
