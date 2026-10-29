package usecase

import (
	"context"
	"fmt"

	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
	"github.com/google/uuid"
)

type UserGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type GetUserByID struct {
	repo UserGetter
}

func NewGetUserByID(repo UserGetter) *GetUserByID {
	return &GetUserByID{
		repo: repo,
	}
}

type GetUserByIDInput struct {
	ID uuid.UUID
}

func (uc *GetUserByID) Execute(ctx context.Context, input GetUserByIDInput) (*domain.User, error) {
	// check uuid id is valid
	id := input.ID

	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
