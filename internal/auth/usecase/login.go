package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/EmotionlessDev/GoEcho/internal/auth/apperror"
	"github.com/EmotionlessDev/GoEcho/internal/helpers/validation"
	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userapperror "github.com/EmotionlessDev/GoEcho/internal/user/apperror"
)

type JWTService interface {
	Create(userID uuid.UUID) (string, error)
}

type PasswordHasher interface {
	Compare(passwordHash, password string) error
}

type UserGetter interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type Login struct {
	users      UserGetter
	password   PasswordHasher
	jwtService JWTService
}

func NewLogin(
	users UserGetter,
	password PasswordHasher,
	jwtService JWTService,
) *Login {
	return &Login{
		users:      users,
		password:   password,
		jwtService: jwtService,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (uc *Login) Execute(
	ctx context.Context,
	input LoginInput,
) (string, error) {
	input.Email = strings.TrimSpace(input.Email)

	if !validation.ValidEmail(input.Email) ||
		!validation.ValidPassword(input.Password) {
		return "", apperror.ErrInvalidInput
	}

	user, err := uc.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, userapperror.ErrUserNotFound) {
			return "", apperror.ErrInvalidCredentials
		}

		return "", fmt.Errorf("get user by email: %w", err)
	}

	if err := uc.password.Compare(user.PasswordHash, input.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", apperror.ErrInvalidCredentials
		}

		return "", fmt.Errorf("compare password: %w", err)
	}

	token, err := uc.jwtService.Create(user.ID)
	if err != nil {
		return "", fmt.Errorf("create JWT token: %w", err)
	}

	return token, nil
}
