package task

import (
	"context"
	"database/sql"
	"errors"
	"web-starter/foundation"
	"web-starter/foundation/appError"

	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/validation"
)

type Task struct {
	Id          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetTaskRequest struct {
	TaskId string `json:"task_id"`
}

func (g *GetTaskRequest) Validation() validation.Validation {
	return validation.NewValidations(
		validation.Value("task_id", g.TaskId, rule.NotBlank[string]()),
	)
}

type GetTaskResponse struct {
	Task *Task `json:"task"`
}

func getTask(ctx context.Context, app *foundation.App, request *GetTaskRequest) (*GetTaskResponse, error) {
	t := &Task{}
	err := app.DB.GetContext(ctx, t, `select id, description from task where id = ?`, request.TaskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.BadRequestError("task not found")
		}
		return nil, appError.InternalServerError(err)
	}
	return &GetTaskResponse{
		Task: t,
	}, nil
}
