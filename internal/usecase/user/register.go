package userusecase

import (
	"context"

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
	repo user.UserRepository
	hasher *hash.BcryptHasher
}

func NewPostgreSQLRegisterUserUseCase(
	repo user.UserRepository,
	hasher *hash.BcryptHasher,
) *PostgreSQLRegisterUserUseCase {
	return &PostgreSQLRegisterUserUseCase{
		repo: repo,
		hasher: hasher,
	}
}

func (usecase *PostgreSQLRegisterUserUseCase) Execute(
	ctx context.Context, input RegisterInput,
) (*user.User, error) {
	// Добавить в миграции проверку длины пароля
	hashPassword, err := usecase.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	userEntity := &user.User{
		Email: input.Email,
		Password: hashPassword,
	}

	err = usecase.repo.Create(ctx, userEntity)
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}
