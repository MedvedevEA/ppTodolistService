package dto

import "github.com/google/uuid"

type AddTaskUser struct {
	TaskId *uuid.UUID `json:"taskId"`
	UserId *uuid.UUID `json:"userId"`
}
type GetTaskUsers struct {
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	TaskId *uuid.UUID `json:"taskId"`
	UserId *uuid.UUID `json:"userId"`
}
