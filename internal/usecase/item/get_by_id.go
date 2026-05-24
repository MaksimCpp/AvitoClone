package itemusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
)

type GetItemByIDUseCase interface {
	Execute(ctx context.Context, id int) (*item.Item, error)
}

type PostgreSQLGetItemByIDUseCase struct {
	repo item.ItemRepository
}

func NewPostgreSQLGetItemByIDUseCase(repo item.ItemRepository) *PostgreSQLGetItemByIDUseCase {
	return &PostgreSQLGetItemByIDUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLGetItemByIDUseCase) Execute(
	ctx context.Context, id int,
) (*item.Item, error) {
	result, err := usecase.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return result, nil
}

