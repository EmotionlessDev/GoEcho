package http

import (
	"context"
	"errors"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	"github.com/EmotionlessDev/GoEcho/internal/auth/apperror"
	"github.com/EmotionlessDev/GoEcho/internal/auth/usecase"
	"github.com/EmotionlessDev/GoEcho/internal/delivery/http/apierror"
)

type Handler struct {
	login *usecase.Login
}

func NewHandler(login *usecase.Login) *Handler {
	return &Handler{
		login: login,
	}
}

func (h *Handler) Login(
	ctx context.Context,
	request httpapi.LoginRequestObject,
) (httpapi.LoginResponseObject, error) {
	if request.Body == nil {
		return httpapi.Login400JSONResponse(httpapi.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Request body is required",
		}), nil
	}

	if request.Body.Email == "" || request.Body.Password == nil {
		return httpapi.Login400JSONResponse(httpapi.ErrorResponse{
			Code:    apierror.InvalidRequest.Code,
			Message: "Email and password are required",
		}), nil
	}

	token, err := h.login.Execute(ctx, usecase.LoginInput{
		Email:    string(request.Body.Email),
		Password: string(*request.Body.Password),
	})
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrInvalidInput):
			return httpapi.Login400JSONResponse(httpapi.ErrorResponse{
				Code:    "INVALID_REQUEST",
				Message: "Invalid email or password format",
			}), nil

		case errors.Is(err, apperror.ErrInvalidCredentials):
			return httpapi.Login401JSONResponse(httpapi.ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			}), nil

		default:
			return nil, err
		}
	}

	// fix the token hardcode
	return httpapi.Login200JSONResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   900,
	}, nil
}
