package apierror

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Write(w http.ResponseWriter, apiErr *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)

	_ = json.NewEncoder(w).Encode(Response{
		Code:    apiErr.Code,
		Message: apiErr.Message,
	})
}
