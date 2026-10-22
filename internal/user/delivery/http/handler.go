package http

import (
	"context"
	"errors"
	"fmt"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	"github.com/EmotionlessDev/GoEcho/internal/user/usecase"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	createUser *usecase.CreateUser
}

func NewHandler(createUser *usecase.CreateUser) *Handler {
	return &Handler{
		createUser: createUser,
	}
}

// Проверка во время компиляции: Handler реализует StrictServerInterface.
var _ httpapi.StrictServerInterface = (*Handler)(nil)

func (h *Handler) CreateUser(
	ctx context.Context,
	request httpapi.CreateUserRequestObject,
) (httpapi.CreateUserResponseObject, error) {
	if request.Body == nil || request.Body.Password == nil {
		return httpapi.CreateUser400JSONResponse{
			Code:    "INVALID_REQUEST",
			Message: "Username, email and password are required",
		}, nil
	}

	user, err := h.createUser.Execute(ctx, usecase.CreateUserInput{
		Username: request.Body.Username,
		Email:    string(request.Body.Email),
		Password: *request.Body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidInput):
			return httpapi.CreateUser400JSONResponse{
				Code:    "INVALID_INPUT",
				Message: "Invalid username, email or password",
			}, nil

		case errors.Is(err, usecase.ErrUserAlreadyExists):
			return httpapi.CreateUser409JSONResponse{
				Code:    "USER_ALREADY_EXISTS",
				Message: "Username or email already exists",
			}, nil

		default:
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	id := openapi_types.UUID(user.ID)
	createdAt := user.CreatedAt

	return httpapi.CreateUser201JSONResponse{
		Id:        &id,
		Username:  user.Username,
		Email:     openapi_types.Email(user.Email),
		CreatedAt: &createdAt,
	}, nil
}

// Пока endpoint не реализован.
func (h *Handler) GetCurrentUser(
	ctx context.Context,
	request httpapi.GetCurrentUserRequestObject,
) (httpapi.GetCurrentUserResponseObject, error) {
	return nil, errors.New("get current user: not implemented")
}

// Пока endpoint не реализован.
func (h *Handler) GetUserByID(
	ctx context.Context,
	request httpapi.GetUserByIDRequestObject,
) (httpapi.GetUserByIDResponseObject, error) {
	return nil, errors.New("get user by ID: not implemented")
}
