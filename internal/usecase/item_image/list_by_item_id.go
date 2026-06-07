package itemimageusecase

import (
	"context"

	itemimage "github.com/MaksimCpp/AvitoClone/internal/domain/item_image"
)

type ListImagesByItemIDUseCase interface {
	Execute(ctx context.Context, itemID int) ([]*itemimage.ItemImage, error)
}

type PostgreSQLListImagesByItemIDUseCase struct {
	repo itemimage.ItemImageRepository
}

func NewPostgreSQLListImagesByItemIDUseCase(
	repo itemimage.ItemImageRepository,
) *PostgreSQLListImagesByItemIDUseCase {
	return &PostgreSQLListImagesByItemIDUseCase{
		repo: repo,
	}
}

func (usecase *PostgreSQLListImagesByItemIDUseCase) Execute(
	ctx context.Context, itemID int,
) ([]*itemimage.ItemImage, error) {
	images, err := usecase.repo.ListByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return images, nil
}
