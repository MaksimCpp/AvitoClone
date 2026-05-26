package itemusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
)

type DeleteItemUseCase interface {
	Execute(ctx context.Context, id int) error
}

type PostgreSQLDeleteItemUseCase struct {
	repo item.ItemRepository
}

func NewPostgreSQLDeleteItemUseCasee(repo item.ItemRepository) *PostgreSQLDeleteItemUseCase {
	return &PostgreSQLDeleteItemUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLDeleteItemUseCase) Execute(
	ctx context.Context, id int,
) error {
	return usecase.repo.Delete(ctx, id)
}

