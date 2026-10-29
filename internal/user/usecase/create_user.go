package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/EmotionlessDev/GoEcho/internal/user/apperror"
	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
)

type UserCreator interface {
	Create(ctx context.Context, user *domain.User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUser struct {
	repo   UserCreator
	hasher PasswordHasher
}

func NewCreateUser(
	repo UserCreator,
	hasher PasswordHasher,
) *CreateUser {
	return &CreateUser{
		repo:   repo,
		hasher: hasher,
	}
}

func (uc *CreateUser) Execute(
	ctx context.Context,
	input CreateUserInput,
) (*domain.User, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)

	if !validUsername(input.Username) ||
		!validEmail(input.Email) ||
		!validPassword(input.Password) {
		return nil, apperror.ErrInvalidInput
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		if errors.Is(err, apperror.ErrUserAlreadyExists) {
			return nil, apperror.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func validUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}

	for _, r := range username {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_') {
			return false
		}
	}

	return true
}

func validEmail(email string) bool {
	return len(email) <= 254 &&
		strings.Contains(email, "@") &&
		!strings.ContainsAny(email, " \t\r\n") &&
		utf8.ValidString(email)
}

func validPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 128
}
