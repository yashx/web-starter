package task

import (
	"encoding/json"
	"net/http"
	"web-starter/foundation"
	"web-starter/foundation/appError"
	"web-starter/foundation/httpHelper"

	"github.com/go-chi/chi/v5"
	"github.com/yashx/shak"
)

type SubRouter struct {
	app *foundation.App
}

func NewSubRouter(app *foundation.App) *SubRouter {
	return &SubRouter{app: app}
}

func (sR SubRouter) BuildHandler() (string, http.Handler) {
	router := chi.NewRouter()
	router.Post("/", sR.getTaskHandlerFn)

	return "/api/task", router
}

func (sR SubRouter) getTaskHandlerFn(writer http.ResponseWriter, httpRequest *http.Request) {
	var request *GetTaskRequest
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		httpHelper.JsonErrorResponse(sR.app, writer, httpRequest, appError.BadRequestError("Invalid Request", err))
		return
	}

	if err := shak.RunValidation(request); err != nil {
		httpHelper.JsonErrorResponse(sR.app, writer, httpRequest, appError.BadRequestErrorFromValidationError(err))
		return
	}

	response, err := getTask(httpRequest.Context(), sR.app, request)
	httpHelper.JsonResponse(sR.app, writer, httpRequest, response, err)
}
