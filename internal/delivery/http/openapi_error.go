package http

import (
	"errors"
	"net/http"

	httpapi "github.com/EmotionlessDev/GoEcho/gen/http"
	"github.com/EmotionlessDev/GoEcho/internal/delivery/http/apierror"
)

func OpenAPIErrorHandler(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	var paramErr *httpapi.InvalidParamFormatError

	if errors.As(err, &paramErr) {
		switch paramErr.ParamName {
		case "userID":
			apierror.Write(w, apierror.InvalidUserID)
			return
		}
	}

	apierror.Write(w, apierror.InvalidRequest)
}
