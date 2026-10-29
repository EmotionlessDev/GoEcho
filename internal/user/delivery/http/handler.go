package http

import (
	"context"
	"errors"
	"fmt"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	"github.com/EmotionlessDev/GoEcho/internal/delivery/http/apierror"
	"github.com/EmotionlessDev/GoEcho/internal/user/apperror"
	"github.com/EmotionlessDev/GoEcho/internal/user/usecase"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	createUser  *usecase.CreateUser
	getUserByID *usecase.GetUserByID
}

func NewHandler(createUser *usecase.CreateUser, getUserByID *usecase.GetUserByID) *Handler {
	return &Handler{
		createUser:  createUser,
		getUserByID: getUserByID,
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
			Code:    apierror.InvalidRequest.Code,
			Message: "Username, email and password are required", // Тут не использую apierror, чтобы добавить детали.
		}, nil
	}

	user, err := h.createUser.Execute(ctx, usecase.CreateUserInput{
		Username: request.Body.Username,
		Email:    string(request.Body.Email),
		Password: *request.Body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrInvalidInput):
			return httpapi.CreateUser400JSONResponse{
				Code:    apierror.InvalidRequest.Code,
				Message: apierror.InvalidRequest.Message,
			}, nil

		case errors.Is(err, apperror.ErrUserAlreadyExists):
			return httpapi.CreateUser409JSONResponse{
				Code:    apierror.UserAlreadyExists.Code,
				Message: apierror.UserAlreadyExists.Message,
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

func (h *Handler) GetUserByID(
	ctx context.Context,
	request httpapi.GetUserByIDRequestObject,
) (httpapi.GetUserByIDResponseObject, error) {
	user, err := h.getUserByID.Execute(ctx, usecase.GetUserByIDInput{
		ID: request.UserID,
	})
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrUserNotFound):
			return httpapi.GetUserByID404JSONResponse{
				Code:    apierror.UserNotFound.Code,
				Message: apierror.UserNotFound.Message,
			}, nil

		default:
			return nil, fmt.Errorf("get user by id: %w", err)
		}
	}

	id := openapi_types.UUID(user.ID)
	createdAt := user.CreatedAt

	return httpapi.GetUserByID200JSONResponse{
		Id:        &id,
		Username:  user.Username,
		CreatedAt: &createdAt,
	}, nil
}
