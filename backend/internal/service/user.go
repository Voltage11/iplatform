package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Voltage11/iplatform/internal/db"
	"github.com/Voltage11/iplatform/internal/domain"
	"github.com/Voltage11/iplatform/internal/types/apperr"
)

// UserRepo — интерфейс сервиса
// Реализация — repo.UserRepo
type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	SoftDelete(ctx context.Context, userID uuid.UUID) error
	SoftUnDelete(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	users      UserRepo
	transactor db.Transactor
	pepper     string
}

func NewUserService(users UserRepo, transactor db.Transactor, pepper string) *UserService {
	return &UserService{
		users:      users,
		transactor: transactor,
		pepper:     pepper,
	}
}

// GetUser — получение пользователя
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

// CreateUser — создание пользователя
func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	// TODO Захешируем пароль

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		existUser, err := s.users.GetByEmail(txCtx, user.Email)
		switch {
		case err == nil && existUser != nil:
			return apperr.NewAlreadyExists(fmt.Sprintf("email занят: %s", existUser.Email), nil)
		case err != nil && !apperr.IsTypeAppError(err, apperr.ErrNotFound):
			return err
		}

		if err := s.users.Create(txCtx, user); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) HashPassword(plain string) (string, error) {
	// pepper дописываем к паролю — так и хранится в хеше bcrypt
	h, err := bcrypt.GenerateFromPassword([]byte(plain+s.pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func (s *UserService) VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain+s.pepper)) == nil
}
