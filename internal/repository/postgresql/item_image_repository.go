package postgresql

import (
	"context"

	itemimage "github.com/MaksimCpp/AvitoClone/internal/domain/item_image"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type ItemImageRepository interface {
// 	Create(ctx context.Context, image *ItemImage) error
// 	ListByItemID(ctx context.Context, itemID int) ([]*ItemImage, error)
// 	Delete(ctx context.Context, id int) error
// }

type PostgreSQLItemImageRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLItemImageRepository(pool *pgxpool.Pool) *PostgreSQLItemImageRepository {
	return &PostgreSQLItemImageRepository{
		pool: pool,
	}
}

func (repo *PostgreSQLItemImageRepository) Create(
	ctx context.Context, image *itemimage.ItemImage,
) error {
	query := `
		INSERT INTO
			item_images (item_id, object_name)
		VALUES
			($1, $2);
	`

	_, err := repo.pool.Exec(
		ctx,
		query,
		image.ItemID,
		image.ObjectName,
	)

	return err
}

func (repo *PostgreSQLItemImageRepository) ListByItemID(
	ctx context.Context, itemID int,
) ([]*itemimage.ItemImage, error) {
	query := `
		SELECT
			id, item_id, object_name
		FROM
			item_images
		WHERE
			item_id = $1;
	`

	var images []*itemimage.ItemImage

	rows, err := repo.pool.Query(
		ctx,
		query,
		itemID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		image := itemimage.ItemImage{}
		err = rows.Scan(
			&image.ID,
			&image.ItemID,
			&image.ObjectName,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}

	return images, nil
}

func (repo *PostgreSQLItemImageRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM item_images
		WHERE id = $1;
	`

	_, err := repo.pool.Exec(
		ctx,
		query,
		id,
	)
	return err
}

func (repo *PostgreSQLItemImageRepository) GetByID(
	ctx context.Context, id int,
) (*itemimage.ItemImage, error) {
	query := `
		SELECT
			id, item_id, object_name
		FROM
			item_images
		WHERE
			id = $1;
	`

	var image itemimage.ItemImage

	err := repo.pool.QueryRow(ctx, query, id).Scan(
		&image.ID,
		&image.ItemID,
		&image.ObjectName,
	)
	if err != nil {
		return nil, err
	}

	return &image, nil
}
