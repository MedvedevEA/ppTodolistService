package dto

import "github.com/google/uuid"

type UpdateStatus struct {
	StatusId *uuid.UUID `json:"statusId"`
	Name     *string    `json:"name"`
}
