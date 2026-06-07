package itemimageusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	itemimage "github.com/MaksimCpp/AvitoClone/internal/domain/item_image"
	"github.com/MaksimCpp/AvitoClone/internal/domain/storage"
)

type DeleteImageUseCase interface {
	Execute(ctx context.Context, imageID int, userID int) error
}

type PostgreSQLDeleteImageUseCase struct {
	imageRepo itemimage.ItemImageRepository
	itemRepo item.ItemRepository
	imageStorage storage.ImageStorage
}

func NewPostgreSQLDeleteImageUseCase(
	imageRepo itemimage.ItemImageRepository,
	itemRepo item.ItemRepository,
	imageStorage storage.ImageStorage,
) *PostgreSQLDeleteImageUseCase {
	return &PostgreSQLDeleteImageUseCase{
		imageRepo: imageRepo,
		itemRepo: itemRepo,
		imageStorage: imageStorage,
	}
}

func (usecase *PostgreSQLDeleteImageUseCase) Execute(
	ctx context.Context, imageID int, userID int,
) error {
	imageEntity, err := usecase.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return err
	}

	itemEntity, err := usecase.itemRepo.GetByID(ctx, imageEntity.ItemID)
	if err != nil {
		return err
	}

	if itemEntity.UserID != userID {
		return item.ErrForbidden
	}

	err = usecase.imageStorage.Delete(ctx, imageEntity.ObjectName)
	if err != nil {
		return err
	}

	return usecase.imageRepo.Delete(ctx, imageID)
}
