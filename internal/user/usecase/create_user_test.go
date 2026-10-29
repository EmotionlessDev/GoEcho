package usecase

import (
	"context"
	"fmt"
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

func TestCreateUser_InvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()

	invalidInputs := []CreateUserInput{
		{Username: "", Email: "", Password: ""},                                // All fields empty
		{Username: "ab", Email: "test@gmail.com", Password: "passwrod123"},     // username too short
		{Username: "testuser", Email: "invalidemail", Password: "password123"}, // invalid email
		{Username: "testuser", Email: "test@gmail.com", Password: "short"},     // password too short
		{Username: "testuser", Email: "", Password: "password123"},             // email empty
		{Username: "testuser", Email: "test@gmail.com", Password: ""},          // password empty
		{Username: "", Email: "test@gmail.com", Password: "password123"},       // username empty
	}

	// Act
	for i, input := range invalidInputs {
		t.Run(fmt.Sprintf("InvalidInputsCase_%d", i), func(t *testing.T) {
			repo := NewMockUserRepository(t)
			hasher := NewMockPasswordHasher(t)

			uc := NewCreateUser(repo, hasher)

			user, err := uc.Execute(ctx, input)

			require.Error(t, err)
			require.Nil(t, user)
			require.Equal(t, ErrInvalidInput, err)
		})
	}
}

func TestCreateUser_UserAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	hasher := NewMockPasswordHasher(t)

	userInput := CreateUserInput{
		Username: "testuser",
		Email:    "test@gmail.com",
		Password: "password123",
	}

	hasher.On("Hash", "password123").Return("hasher_password", nil).Once()

	repo.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == "testuser" && user.Email == "test@gmail.com" && user.PasswordHash == "hasher_password"
	})).Return(ErrUserAlreadyExists).Once()

	uc := NewCreateUser(repo, hasher)

	// Act
	user, err := uc.Execute(ctx, userInput)

	// Aseert
	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, ErrUserAlreadyExists, err)

	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestCreateUser_HashError(t *testing.T) {
	// Arrange
	ctx := context.Background()

	repo := NewMockUserRepository(t)
	hasher := NewMockPasswordHasher(t)

	userInput := CreateUserInput{
		Username: "testuser",
		Email:    "test@gmail.com",
		Password: "password123",
	}

	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	hasher.On("Hash", "password123").Return("", fmt.Errorf("hash error")).Once()

	uc := NewCreateUser(repo, hasher)

	// Act
	user, err := uc.Execute(ctx, userInput)

	// Assert
	require.Error(t, err)
	require.Nil(t, user)
	require.Contains(t, err.Error(), "hash password")

	repo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}
