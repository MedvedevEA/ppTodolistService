package repository

import (
	"ppTodolistService/internal/entity"
	repoDto "ppTodolistService/internal/repository/dto"

	"github.com/google/uuid"
)

type Repository interface {
	AddMessage(dto *repoDto.AddMessage) (*entity.Message, error)
	GetMessage(messageId *uuid.UUID) (*entity.Message, error)
	GetMessages(dto *repoDto.GetMessages) ([]*entity.Message, error)
	UpdateMessage(dto *repoDto.UpdateMessage) (*entity.Message, error)
	RemoveMessage(messageId *uuid.UUID) error

	AddStatus(name string) (*entity.Status, error)
	GetStatus(statusId *uuid.UUID) (*entity.Status, error)
	GetStatuses() ([]*entity.Status, error)
	UpdateStatus(dto *repoDto.UpdateStatus) (*entity.Status, error)
	RemoveStatus(statusId *uuid.UUID) error

	AddTask(dto *repoDto.AddTask) (*entity.Task, error)
	GetTask(taskId *uuid.UUID) (*entity.Task, error)
	GetTasks(dto *repoDto.GetTasks) ([]*entity.Task, error)
	UpdateTask(dto *repoDto.UpdateTask) (*entity.Task, error)
	RemoveTask(taskId *uuid.UUID) error

	AddTaskUser(dto *repoDto.AddTaskUser) (*entity.TaskUser, error)
	GetTaskUsers(dto *repoDto.GetTaskUsers) ([]*entity.TaskUser, error)
	RemoveTaskUser(taskUserID *uuid.UUID) error

	AddUserWithUserId(dto *repoDto.AddUser) (*entity.User, error)
	GetUsers(dto *repoDto.GetUsers) ([]*entity.User, error)
	RemoveUser(userId *uuid.UUID) error
}
