package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *User) error
}

type pgStorage struct{}

func NewStorage() *pgStorage {
	return &pgStorage{}
}

const createUserQuery = `
	INSERT INTO users (username, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
`

func (s *pgStorage) CreateUser(ctx context.Context, tx pgx.Tx, user *User) error {
	err := tx.QueryRow(
		ctx,
		createUserQuery,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
