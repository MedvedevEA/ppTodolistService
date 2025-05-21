package store

import (
	"context"
	"database/sql"
	"log/slog"
	"ppTodolistService/internal/entity"
	"ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	addStatusQuery = `
INSERT INTO status (name) 
VALUES ($1) RETURNING *;`
	getStatusQuery = `
SELECT * FROM status WHERE status_id=$1;`
	getStatusesQuery = `
SELECT * FROM status;`
	updateStatusQuery = `
UPDATE status 
SET name = CASE WHEN $2::character varying IS NULL THEN name ELSE $2 END
WHERE status_id=$1
RETURNING *;`
	removeStatusQuery = `
DELETE FROM status 
WHERE status_id=$1;`
)

func (s *Store) AddStatus(name string) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), addStatusQuery, name).Scan(&status.StatusId, &status.Name)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddStatus"))
		return nil, repoErr.ErrInternalServerError
	}
	return status, nil
}
func (s *Store) GetStatus(statusId *uuid.UUID) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), getStatusQuery, statusId).Scan(&status.StatusId, &status.Name)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetStatus"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return status, nil
}
func (s *Store) GetStatuses() ([]*entity.Status, error) {
	rows, err := s.pool.Query(context.Background(), getStatusesQuery)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetStatuses"))
		return nil, repoErr.ErrInternalServerError
	}
	defer rows.Close()
	statuses := make([]*entity.Status, 0)
	for rows.Next() {
		status := new(entity.Status)
		err := rows.Scan(&status.StatusId, &status.Name)
		if err != nil {
			s.lg.Error(err.Error(), slog.String("owner", "store.GetStatuses"))
			return nil, repoErr.ErrInternalServerError
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
func (s *Store) UpdateStatus(dto *dto.UpdateStatus) (*entity.Status, error) {
	status := new(entity.Status)
	err := s.pool.QueryRow(context.Background(), updateStatusQuery, dto.StatusId, dto.Name).Scan(&status.StatusId, &status.Name)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.UpdateStatus"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return status, nil
}
func (s *Store) RemoveStatus(statusId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeStatusQuery, statusId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveStatus"))
		return repoErr.ErrInternalServerError
	}
	if result.RowsAffected() == 0 {
		s.lg.Error(sql.ErrNoRows.Error(), slog.String("owner", "store.RemoveStatus"))
		return repoErr.ErrRecordNotFound
	}
	return nil
}
