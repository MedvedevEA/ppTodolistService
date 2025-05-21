package dto

import "github.com/google/uuid"

type AddMessage struct {
	TaskId *uuid.UUID `json:"taskId"`
	UserId *uuid.UUID `json:"userId"`
	Text   string     `json:"text"`
}
type GetMessages struct {
	Offset int        `json:"offset"`
	Limit  int        `json:"limit"`
	TaskId *uuid.UUID `json:"taskId"`
}
type UpdateMessage struct {
	MessageId *uuid.UUID `json:"messageId"`
	Text      *string    `json:"text"`
}
