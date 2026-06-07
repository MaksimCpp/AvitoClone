package itemimageusecase

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/MaksimCpp/AvitoClone/internal/config"
	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	itemimage "github.com/MaksimCpp/AvitoClone/internal/domain/item_image"
	"github.com/MaksimCpp/AvitoClone/internal/domain/storage"
	"github.com/google/uuid"
)

type UploadInput struct {
	ItemID int
	UserID int
	File multipart.File
	Filename string
	ContentType string
	Size int64
}

type UploadImageUseCase interface {
	Execute(ctx context.Context, input UploadInput) (string, error)
}

type PostgreSQLUploadImageUseCase struct {
	imageRepo itemimage.ItemImageRepository
	itemRepo item.ItemRepository
	imageStorage storage.ImageStorage
	cfg *config.Config
}

func NewPostgreSQLUploadImageUseCase(
	imageRepo itemimage.ItemImageRepository,
	itemRepo item.ItemRepository,
	imageStorage storage.ImageStorage,
	cfg *config.Config,
) *PostgreSQLUploadImageUseCase {
	return &PostgreSQLUploadImageUseCase{
		imageRepo: imageRepo,
		itemRepo: itemRepo,
		imageStorage: imageStorage,
		cfg: cfg,
	}
}

func (usecase *PostgreSQLUploadImageUseCase) Execute(
	ctx context.Context, input UploadInput,
) (string, error) {
	itemEntity, err := usecase.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		return "", err
	}

	if itemEntity.UserID != input.UserID {
		return "", itemimage.ErrPhotoNotOwned
	}

	objectName := fmt.Sprintf(
		"%s-%s",
		uuid.New().String(),
		input.Filename,
	)

	url, err := usecase.imageStorage.Upload(
		ctx,
		input.File,
		objectName,
		input.ContentType,
		input.Size,
		usecase.cfg.MinioEndpoint,
	)
	if err != nil {
		return "", err
	}

	imageEntity := itemimage.ItemImage{
		ItemID: input.ItemID,
		ObjectName: objectName,
	}

	err = usecase.imageRepo.Create(ctx, &imageEntity)
	if err != nil {
		return "", err
	}

	return url, nil
}
