package user

import (
	"context"
	"go-gorm/internal/database/model"
	"go-gorm/internal/database/query"
)

type RepositoryUser interface {
	Create(ctx context.Context, user *model.User) error
	CreateBatch(ctx context.Context, users []*model.User, batchSize int) error
}

type userRepository struct {
	query *query.Query
}

func NewUserRepository(q *query.Query) *userRepository {
	return &userRepository{query: q}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := r.query.User.WithContext(ctx).Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) CreateBatch(ctx context.Context, users []*model.User, batchSize int) error {
	err := r.query.User.WithContext(ctx).CreateInBatches(users, batchSize)
	if err != nil {
		return err
	}
	return nil
}
