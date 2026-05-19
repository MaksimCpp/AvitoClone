package userusecase

import (
	"context"
	"errors"

	"github.com/MaksimCpp/AvitoClone/internal/domain/user"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/hash"
)

type RegisterInput struct {
	Email string
	Password string
}

type RegisterUserUseCase interface {
	Execute(ctx context.Context, input RegisterInput) error
}

type PostgreSQLRegisterUserUseCase struct {
	repo *user.UserRepository
	hasher *hash.BcryptHasher
}

func NewPostgreSQLRegisterUserUseCase(
	repo *user.UserRepository,
	hasher *hash.BcryptHasher,
) *PostgreSQLRegisterUserUseCase {
	return &PostgreSQLRegisterUserUseCase{
		repo: repo,
		hasher: hasher,
	}
}

func (usecase *PostgreSQLRegisterUserUseCase) Execute(ctx context.Context, input RegisterInput) error {
	return errors.New("error")
}
