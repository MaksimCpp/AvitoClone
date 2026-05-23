package userusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/user"
)

type GetMeUseCase interface {
	Execute(ctx context.Context, id int) (*user.User, error)
}

type PostgreSQLGetMeUseCase struct {
	repo user.UserRepository
}

func NewPostgreSQLGetMeUseCase(
	repo user.UserRepository,
) *PostgreSQLGetMeUseCase {
	return &PostgreSQLGetMeUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLGetMeUseCase) Execute(ctx context.Context, id int) (*user.User, error) {
	userEntity, err := usecase.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return userEntity, nil
}
