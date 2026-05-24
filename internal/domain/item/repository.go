package item

import "context"

type ItemRepository interface {
	Create(ctx context.Context, itemEntity *Item) error
	GetByID(ctx context.Context, id int) (*Item, error)
	List(ctx context.Context, limit int, offset int) ([]*Item, error)
	ListByUserID(ctx context.Context, userID int) ([]*Item, error)
	Delete(ctx context.Context, id int) error
}
