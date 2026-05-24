package itemusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
)

type CreateItemUseCase interface {
	Execute(ctx context.Context, input *item.Item) (*item.Item, error)
}

type PostgreSQLCreateItemUseCase struct {
	repo item.ItemRepository
}

func NewPostgreSQLCreateItemUseCase(repo item.ItemRepository) *PostgreSQLCreateItemUseCase {
	return &PostgreSQLCreateItemUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLCreateItemUseCase) Execute(
	ctx context.Context, input *item.Item,
) (*item.Item, error) {
	err := usecase.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

