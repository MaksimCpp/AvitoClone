package item

import (
	"time"

	itemimage "github.com/MaksimCpp/AvitoClone/internal/domain/item_image"
)

type Item struct {
	ID int
	UserID int
	Title string
	Description string
	Price float64
	CreatedAt time.Time
	Images []itemimage.ItemImage
}
