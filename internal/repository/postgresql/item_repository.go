package postgresql

import (
	"context"
	"errors"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	"github.com/jackc/pgx/v5"
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
	).Scan(
		&itemEntity.ID,
		&itemEntity.CreatedAt,
	)

}

func (repo *PostgreSQLItemRepository) GetByID(ctx context.Context, id int) (*item.Item, error) {
	query := `
		SELECT
			id, user_id, title, description, price
		FROM
			items
		WHERE id = $1;
	`

	var itemEntity item.Item

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&itemEntity.ID,
		&itemEntity.UserID,
		&itemEntity.Title,
		&itemEntity.Description,
		&itemEntity.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, item.ErrItemNotFound
		}
		return nil, err
	}

	return &itemEntity, nil
}

func (repo *PostgreSQLItemRepository) List(
	ctx context.Context, limit int, offset int,
) ([]*item.Item, error) {
	query := `
		SELECT
			title, price
		FROM
			items
		LIMIT $1 OFFSET $2;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*item.Item

	for rows.Next() {
		itemEntity := item.Item{}
		rows.Scan(
			&itemEntity.Title,
			&itemEntity.Price,
		)

		items = append(items, &itemEntity)
	}

	return items, nil
}

func (repo *PostgreSQLItemRepository) ListByUserID(
	ctx context.Context, userID int,
) ([]*item.Item, error) {
	query := `
		SELECT
			title, price
		FROM
			items
		WHERE
			user_id = $1;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*item.Item

	for rows.Next() {
		itemEntity := item.Item{}
		rows.Scan(
			&itemEntity.Title,
			&itemEntity.Price,
		)

		items = append(items, &itemEntity)
	}

	return items, nil
}

func (repo *PostgreSQLItemRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM items
		WHERE id = $1;
	`

	result, err := repo.pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return item.ErrItemNotFound
	}

	return nil
}
