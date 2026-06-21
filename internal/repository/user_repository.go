package repository

import (
	"acc_backend/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrUserAlreadyExists = errors.New("user with this email already exists")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(fullName string, email string, password string) (string, error) {
	user := model.User{FullName: fullName, Email: email, Password: password}

	ctx := context.Background()
	err := gorm.G[model.User](r.db).Create(ctx, &user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", ErrUserAlreadyExists
		}
		return "", err
	}
	return user.ID, nil
}

func (r *UserRepository) GetById(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
