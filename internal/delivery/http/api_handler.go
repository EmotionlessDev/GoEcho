package http

import (
	"context"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	userHttp "github.com/EmotionlessDev/GoEcho/internal/user/delivery/http"
)

type APIHandler struct {
	users *userHttp.Handler
}

func NewAPIHandler(users *userHttp.Handler) *APIHandler {
	return &APIHandler{
		users: users,
	}
}

var _ httpapi.StrictServerInterface = (*APIHandler)(nil)

func (h *APIHandler) CreateUser(
	ctx context.Context,
	request httpapi.CreateUserRequestObject,
) (httpapi.CreateUserResponseObject, error) {
	return h.users.CreateUser(ctx, request)
}

func (h *APIHandler) GetCurrentUser(
	ctx context.Context,
	request httpapi.GetCurrentUserRequestObject,
) (httpapi.GetCurrentUserResponseObject, error) {
	return h.users.GetCurrentUser(ctx, request)
}

func (h *APIHandler) GetUserByID(
	ctx context.Context,
	request httpapi.GetUserByIDRequestObject,
) (httpapi.GetUserByIDResponseObject, error) {
	return h.users.GetUserByID(ctx, request)
}
