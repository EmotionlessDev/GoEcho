package usecase

import (
	"context"
	"testing"

	"github.com/EmotionlessDev/GoEcho/internal/user/domain"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	hasher := NewMockPasswordHasher(t)

	hasher.On("Hash", "password123").Return("hashed_password", nil).Once()

	repo.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == "testuser" && user.Email == "test@gmail.com" && user.PasswordHash == "hashed_password"
	})).Return(nil).Once()

	userInput := CreateUserInput{
		Username: "testuser",
		Email:    "test@gmail.com",
		Password: "password123",
	}

	uc := NewCreateUser(repo, hasher)

	// Act
	user, err := uc.Execute(ctx, userInput)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, userInput.Username, user.Username)
	require.Equal(t, userInput.Email, user.Email)
	require.Equal(t, "hashed_password", user.PasswordHash)

	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}
