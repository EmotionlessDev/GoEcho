package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/GoEcho/internal/user/apperror"
	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

const createUserQuery = `
    INSERT INTO users (username, email, password_hash)
    VALUES ($1, $2, $3)
    RETURNING id, created_at, updated_at
`

func (r *PostgresUserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	err := r.pool.QueryRow(
		ctx,
		createUserQuery,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

const getUserByIDQuery = `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1
`

func (r *PostgresUserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, getUserByIDQuery, id)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}
