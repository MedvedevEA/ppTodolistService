package dto

import "github.com/google/uuid"

type AddUser struct {
	UserId *uuid.UUID `json:"user_id"`
	Name   string     `json:"name"`
}
type GetUsers struct {
	Offset int     `query:"offset" validate:"gte=0"`
	Limit  int     `query:"limit" validate:"gte=0"`
	Name   *string `query:"name" validate:"omitempty"`
}
type RemoveUser struct {
	UserId *uuid.UUID `uri:"user_id"`
}
