package store

import (
	"context"
	"database/sql"
	"log/slog"
	"ppTodolistService/internal/entity"
	repoDto "ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	addTaskUserQuery = `
INSERT INTO task_user (task_id, user_id) 
VALUES ($1, $2) RETURNING *;`
	getTaskUsersQuery = `
SELECT * 
FROM task_user 
WHERE ($3::uuid IS NULL OR task_id=$3) AND ($4::uuid IS null OR user_id=$4) 
OFFSET $1 LIMIT $2;`
	removeTaskUserQuery = `
DELETE FROM task_user
WHERE task_user_id=$1;`
)

func (s *Store) AddTaskUser(dto *repoDto.AddTaskUser) (*entity.TaskUser, error) {
	taskUser := new(entity.TaskUser)
	err := s.pool.QueryRow(context.Background(), addTaskUserQuery, dto.TaskId, dto.UserId).Scan(&taskUser.TaskUserId, &taskUser.UserId, &taskUser.TaskId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddTaskUsers"))
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == "23505" {
			return nil, repoErr.ErrUniqueViolation
		}
		return nil, repoErr.ErrInternalServerError

	}
	return taskUser, nil
}
func (s *Store) GetTaskUsers(dto *repoDto.GetTaskUsers) ([]*entity.TaskUser, error) {
	rows, err := s.pool.Query(context.Background(), getTaskUsersQuery, dto.Offset, dto.Limit, dto.TaskId, dto.UserId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetTaskUsers"))
		return nil, repoErr.ErrInternalServerError
	}
	defer rows.Close()
	taskUsers := make([]*entity.TaskUser, 0)
	for rows.Next() {
		taskUser := new(entity.TaskUser)
		err := rows.Scan(&taskUser.TaskUserId, &taskUser.TaskId, &taskUser.UserId)
		if err != nil {
			s.lg.Error(err.Error(), slog.String("owner", "store.GetTaskUsers"))
			return nil, repoErr.ErrInternalServerError
		}
		taskUsers = append(taskUsers, taskUser)
	}
	return taskUsers, nil
}

func (s *Store) RemoveTaskUser(TaskUserId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeTaskUserQuery, TaskUserId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveTaskUser"))
		return repoErr.ErrInternalServerError
	}
	if result.RowsAffected() == 0 {
		s.lg.Error(sql.ErrNoRows.Error(), slog.String("owner", "store.RemoveTaskUser"))
		return repoErr.ErrRecordNotFound
	}
	return nil
}
