package itemusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
)

type ListItemsByUserIDUseCase interface {
	Execute(ctx context.Context, userID int) ([]*item.Item, error)
}

type PostgreSQLListItemsByUserIDUseCase struct {
	repo item.ItemRepository
}

func NewPostgreSQLListItemsByUserIDUseCase(repo item.ItemRepository) *PostgreSQLListItemsByUserIDUseCase {
	return &PostgreSQLListItemsByUserIDUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLListItemsByUserIDUseCase) Execute(
	ctx context.Context, userID int,
) ([]*item.Item, error) {
	result, err := usecase.repo.ListByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return result, nil
}

