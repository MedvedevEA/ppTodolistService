package store

import (
	"context"
	"database/sql"
	"log/slog"
	"ppTodolistService/internal/entity"
	repoDto "ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	addTaskQuery = `
INSERT INTO task (status_id, title, description) 
VALUES ($1, $2, $3) 
RETURNING *;`
	getTaskQuery = `
SELECT * FROM task WHERE task_id=$1;`
	getTasksQuery = `
SELECT * FROM task 
WHERE $3::uuid IS null OR status_id = $3 
OFFSET $1 LIMIT $2;`
	updateTaskQuery = `
UPDATE task 
SET 
status_id = CASE WHEN $2::uuid IS NULL THEN status_id ELSE $2 END,
title = CASE WHEN $3::character varying IS NULL THEN title ELSE $3 END,
description = CASE WHEN $4::character varying IS NULL THEN description ELSE $4 END
WHERE task_id=$1
RETURNING *;`
	removeTaskQuery = `
DELETE FROM task 
WHERE task_id=$1;`
)

func (s *Store) AddTask(dto *repoDto.AddTask) (*entity.Task, error) {
	task := new(entity.Task)
	err := s.pool.QueryRow(context.Background(), addTaskQuery, dto.StatusId, dto.Title, dto.Description).Scan(&task.TaskId, &task.StatusId, &task.Title, &task.Description)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddTask"))
		return nil, repoErr.ErrInternalServerError
	}
	return task, nil
}
func (s *Store) GetTask(taskId *uuid.UUID) (*entity.Task, error) {
	task := new(entity.Task)
	err := s.pool.QueryRow(context.Background(), getTaskQuery, taskId).Scan(&task.TaskId, &task.StatusId, &task.Title, &task.Description)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetTask"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return task, nil
}
func (s *Store) GetTasks(dto *repoDto.GetTasks) ([]*entity.Task, error) {
	rows, err := s.pool.Query(context.Background(), getTasksQuery, dto.Offset, dto.Limit, dto.StatusId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetTasks"))
		return nil, repoErr.ErrInternalServerError
	}
	defer rows.Close()
	tasks := make([]*entity.Task, 0)
	for rows.Next() {
		task := new(entity.Task)
		err := rows.Scan(&task.TaskId, &task.StatusId, &task.Title, &task.Description)
		if err != nil {
			s.lg.Error(err.Error(), slog.String("owner", "store.GetTasks"))
			return nil, repoErr.ErrInternalServerError
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (s *Store) UpdateTask(dto *repoDto.UpdateTask) (*entity.Task, error) {
	task := new(entity.Task)
	err := s.pool.QueryRow(context.Background(), updateTaskQuery, dto.TaskId, dto.StatusId, dto.Title, dto.Description).Scan(&task.TaskId, &task.StatusId, &task.Title, &task.Description)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.UpdateTask"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return task, nil
}

func (s *Store) RemoveTask(taskId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeTaskQuery, taskId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveTask"))
		return repoErr.ErrInternalServerError
	}
	if result.RowsAffected() == 0 {
		s.lg.Error(sql.ErrNoRows.Error(), slog.String("owner", "store.RemoveTask"))
		return repoErr.ErrRecordNotFound
	}

	return nil

}
