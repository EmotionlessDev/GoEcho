package http

import (
	"context"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	authHttp "github.com/EmotionlessDev/GoEcho/internal/auth/delivery/http"
	userHttp "github.com/EmotionlessDev/GoEcho/internal/user/delivery/http"
)

type APIHandler struct {
	users *userHttp.Handler
	auth  *authHttp.Handler
}

func NewAPIHandler(users *userHttp.Handler, auth *authHttp.Handler) *APIHandler {
	return &APIHandler{
		users: users,
		auth:  auth,
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

func (h *APIHandler) Login(
	ctx context.Context,
	request httpapi.LoginRequestObject,
) (httpapi.LoginResponseObject, error) {
	return h.auth.Login(ctx, request)
}
