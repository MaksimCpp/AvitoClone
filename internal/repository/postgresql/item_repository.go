package postgresql

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLItemRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLItemRepository(pool *pgxpool.Pool) *PostgreSQLItemRepository {
	return &PostgreSQLItemRepository{
		pool: pool,
	}
}

// Create(ctx context.Context, itemEntity *Item) error
// GetByID(ctx context.Context, id int) (*Item, error)
// List(ctx context.Context, limit int, offset int) ([]*Item, error)
// Delete(ctx context.Context, id int) error

func (repo *PostgreSQLItemRepository) Create(ctx context.Context, itemEntity *item.Item) error {
	query := `
		INSERT INTO
			items (user_id, title, description, price)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id, created_at;
	`

	return repo.pool.QueryRow(
		ctx,
		query,
		itemEntity.UserID,
		itemEntity.Title,
		itemEntity.Description,
		itemEntity.Price,
		itemEntity.Price,
	).Scan(
		&itemEntity.ID,
		&itemEntity.CreatedAt,
	)

}
