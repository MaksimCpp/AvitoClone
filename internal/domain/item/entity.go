package item

import "time"

type Item struct {
	ID int
	UserID int
	Title string
	Description string
	Price float64
	CreatedAt time.Time
	Images []ItemImage
}

type ItemImage struct {
	ID int
	ItemID int
	ImageURL string
}
