package user

import (
	"context"
	"errors"
	"go-gorm/internal/database/model"
	"go-gorm/internal/database/query"

	"gorm.io/gorm"
)

type RepositoryUser interface {
	Create(ctx context.Context, user *model.User) error
	CreateBatch(ctx context.Context, users []*model.User, batchSize int) error
	Transaction(ctx context.Context, user *model.User) error
	ManualTransaction(ctx context.Context, user *model.User) error
	GetFirstID(ctx context.Context, id int64) (*model.User, error)
	GetFirstEmail(ctx context.Context) (*model.User, error)
	GetLastEmail(ctx context.Context) (*model.User, error)
	GetAllFind(ctx context.Context) ([]*model.User, error)
	GetAllPaginate(ctx context.Context, limit int, offset int) ([]*model.User, error)
	SearchUser(ctx context.Context, name string, email string) ([]*model.User, error)
	SearcUserOrEmail(ctx context.Context, name string, email string) (*model.User, error)
	SearchNotName(ctx context.Context, name string) ([]*model.User, error)
	SelectFieldName(ctx context.Context) ([]*model.User, error)
	StructCondition(ctx context.Context, name string, email string) (*model.User, error)
	PaginationOrderName(ctx context.Context, limit int, offset int) ([]*model.User, error)
	Update(ctx context.Context, id int64, name string) error 
}

type userRepository struct {
	db    *gorm.DB
	query *query.Query
}

func NewUserRepository(q *query.Query, db *gorm.DB) *userRepository {
	return &userRepository{query: q, db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.query.User.WithContext(ctx).Create(user)
}

func (r *userRepository) CreateBatch(ctx context.Context, users []*model.User, batchSize int) error {
	return r.query.User.WithContext(ctx).CreateInBatches(users, batchSize)

}

func (r *userRepository) Transaction(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		qtx := query.Use(tx)
		return qtx.User.WithContext(ctx).Create(user)
	})
}

func (r *userRepository) ManualTransaction(ctx context.Context, user *model.User) error {
	tx := r.db.WithContext(ctx).Begin()
	qtx := query.Use(tx)
	err := qtx.User.WithContext(ctx).Create(user)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// First() -> digunakan untuk mengambil data pertama yang diurutkan berdasarkan id -> ASC

func (r *userRepository) GetFirstID(ctx context.Context, id int64) (*model.User, error) {
	result, err := r.query.User.WithContext(ctx).Where(r.query.User.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetFirstEmail(ctx context.Context) (*model.User, error) {
	result, err := r.query.User.WithContext(ctx).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// Take() -> untuk get satu data aja spesifik
func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.query.User.WithContext(ctx).Where(r.query.User.ID.Eq(id)).Take()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Last () -> kebalikannya first DESC ID -> untuk category atau jenis yang banyak data
func (r *userRepository) GetLastEmail(ctx context.Context) (*model.User, error) {
	user, err := r.query.User.WithContext(ctx).Last()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAllFind(ctx context.Context) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	return q.Find()
}

func (r *userRepository) GetAllPaginate(ctx context.Context, limit int, offset int) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx).Limit(limit).Offset(offset)
	return q.Find()
}

// Advanced Query
func (r *userRepository) SearchUser(ctx context.Context, name string, email string) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	if name != "" {
		q = q.Where(u.Name.Like("%" + name + "%")).Where(u.Email.Eq(email))
	}
	return q.Find()
}

// OR
func (r *userRepository) SearcUserOrEmail(ctx context.Context, name string, email string) (*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	if name != "" && email != "" {
		q = q.Where(u.Name.Eq(name)).Or(u.Email.Eq(email))
	}
	return q.Take()
}
func (r *userRepository) SearchNotName(ctx context.Context, name string) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	if name != "" {
		q = q.Not(u.Name.Eq(name))
	}
	return q.Find()
}

func (r *userRepository) SelectFieldName(ctx context.Context) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx).Select(u.Name)
	return q.Find()
}

func (r *userRepository) StructCondition(ctx context.Context, name string, email string) (*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	if name != "" && email != "" {
		q = q.Where(
			u.Name.Eq(name),
			u.Email.Eq(email),
		)
	}
	return q.First()
}

func (r *userRepository) PaginationOrderName(ctx context.Context, limit int, offset int) ([]*model.User, error) {
	u := r.query.User
	q := u.WithContext(ctx)
	q = q.Order(u.ID.Asc()).Limit(limit).Offset(offset)
	return q.Find()
}

func (r *userRepository) Update(ctx context.Context, id int64, name string) error {
	u := r.query.User
	q := u.WithContext(ctx).Where(u.ID.Eq(id))
	if name != "" {
		_, err := q.Update(u.Name, name)
		if err != nil {
			return err
		}
	}
	return nil
}
