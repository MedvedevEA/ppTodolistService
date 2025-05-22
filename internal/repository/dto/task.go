package dto

import "github.com/google/uuid"

type AddTask struct {
	UserId      *uuid.UUID `json:"userId"`
	StatusId    *uuid.UUID `json:"statusId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}
type GetTasks struct {
	Offset   int        `json:"offset"`
	Limit    int        `json:"limit"`
	StatusId *uuid.UUID `json:"statusId"`
}
type UpdateTask struct {
	TaskId      *uuid.UUID `json:"taskId"`
	StatusId    *uuid.UUID `json:"statusId"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
}
