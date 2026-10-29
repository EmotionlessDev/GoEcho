package apierror

import "net/http"

var (
	InvalidRequest = &Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_REQUEST",
		Message: "Invalid request",
	}

	InvalidUserID = &Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_USER_ID",
		Message: "User ID must be a valid UUID",
	}

	UserNotFound = &Error{
		Status:  http.StatusNotFound,
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
	}

	UserAlreadyExists = &Error{
		Status:  http.StatusConflict,
		Code:    "USER_ALREADY_EXISTS",
		Message: "Username or email already exists",
	}

	Internal = &Error{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
	}
)
