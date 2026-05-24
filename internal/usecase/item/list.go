package itemusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
)

type ListItemsUseCase interface {
	Execute(ctx context.Context, limit int, offset int) ([]*item.Item, error)
}

type PostgreSQLListItemsUseCase struct {
	repo item.ItemRepository
}

func NewPostgreSQLListItemsUseCase(repo item.ItemRepository) *PostgreSQLListItemsUseCase {
	return &PostgreSQLListItemsUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLListItemsUseCase) Execute(
	ctx context.Context, limit int, offset int,
) ([]*item.Item, error) {
	result, err := usecase.repo.List(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	return result, nil
}

