package httpHelper

import (
	"errors"
	"net/http"
	"web-starter/foundation"
	"web-starter/foundation/appError"

	"github.com/go-chi/render"
)

func JsonResponse(app *foundation.App, responseWriter http.ResponseWriter, httpRequest *http.Request, response any, err error) {
	if err != nil {
		JsonErrorResponse(app, responseWriter, httpRequest, err)
	} else {
		render.JSON(responseWriter, httpRequest, response)
	}
}

func JsonErrorResponse(app *foundation.App, responseWriter http.ResponseWriter, httpRequest *http.Request, err error) {
	var aErr *appError.AppError
	if !errors.As(err, &aErr) {
		aErr = appError.InternalServerError()
	}
	render.Status(httpRequest, aErr.HttpStatus)
	render.JSON(responseWriter, httpRequest, aErr)
}
