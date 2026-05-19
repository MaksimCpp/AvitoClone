package postgresql

import (
	"context"
	"errors"

	"github.com/MaksimCpp/AvitoClone/internal/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLUserRepository(pool *pgxpool.Pool) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{
		pool: pool,
	}
}

func (repo *PostgreSQLUserRepository) Create(ctx context.Context, userEntity *user.User) error {
	query := `
		INSERT INTO
			users (email, password)
		VALUES
			($1, $2)
		RETURNING
			id, email, created_at;
	`
	err := repo.pool.QueryRow(
		ctx,
		query,
		userEntity.Email,
		userEntity.Password,
	).Scan(
		&userEntity.ID,
		&userEntity.Email,
		&userEntity.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		uniqueViolationCode := "23505"

		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return user.ErrUserAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (repo *PostgreSQLUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT
			id, email, password, created_at
		FROM
			users
		WHERE
			email = $1;
	`
	var userEntity user.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&userEntity.ID,
		&userEntity.Email,
		&userEntity.Password,
		&userEntity.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &userEntity, nil
}

func (repo *PostgreSQLUserRepository) GetById(ctx context.Context, id int) (*user.User, error) {
	query := `
		SELECT
			id, email, password, created_at
		FROM
			users
		WHERE
			id = $1;
	`
	var userEntity user.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&userEntity.ID,
		&userEntity.Email,
		&userEntity.Password,
		&userEntity.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &userEntity, nil
}