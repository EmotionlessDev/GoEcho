package repository

import (
	"context"
	"fmt"

	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
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
