package httpHelper

import (
	"net/http"
	"web-starter/foundation/appError"
)

func MustQueryParam(r *http.Request, param string) (string, error) {
	paramStrList, ok := r.URL.Query()[param]
	if !ok {
		return "", appError.BadRequestError(param + " query param required")
	}

	return paramStrList[0], nil
}
