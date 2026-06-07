package itemimage

import "context"

type ItemImageRepository interface {
	Create(ctx context.Context, image *ItemImage) error
	ListByItemID(ctx context.Context, itemID int) ([]*ItemImage, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*ItemImage, error)
}

