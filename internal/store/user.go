package store

import (
	"context"
	"fmt"
	"ppTodolistService/internal/entity"
	repoDto "ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"

	"github.com/google/uuid"
)

const (
	addUserQuery = `
INSERT INTO "user" (user_id,name) 
VALUES ($1,$2) 
RETURNING *;`
	getUsersQuery = `
SELECT * FROM "user" 
WHERE $3::character varying IS NULL OR name ILIKE '%'||$3||'%'
OFFSET $1 LIMIT $2;`
	removeUserQuery = `
DELETE FROM "user" 
WHERE user_id=$1;`
)

func (s *Store) AddUserWithUserId(dto *repoDto.AddUser) (*entity.User, error) {
	user := new(entity.User)
	err := s.pool.QueryRow(context.Background(), addUserQuery, dto.UserId, dto.Name).Scan(&user.UserId, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("store: AddUser: error: %w: %v", repoErr.ErrInternalServerError, err)
	}
	return user, nil
}
func (s *Store) GetUsers(dto *repoDto.GetUsers) ([]*entity.User, error) {
	rows, err := s.pool.Query(context.Background(), getUsersQuery, dto.Offset, dto.Limit, dto.Name)
	if err != nil {
		return nil, fmt.Errorf("store: GetUsers: error: %w: %v", repoErr.ErrInternalServerError, err)
	}
	defer rows.Close()
	users := make([]*entity.User, 0)
	for rows.Next() {
		user := new(entity.User)
		err := rows.Scan(&user.UserId, &user.Name)
		if err != nil {
			return nil, fmt.Errorf("store: GetUsers: error: %w: %v", repoErr.ErrInternalServerError, err)
		}
		users = append(users, user)
	}
	return users, nil
}
func (s *Store) RemoveUser(userId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeUserQuery, userId)
	if err != nil {
		return fmt.Errorf("store: RemoveUser: error: %w: %v", repoErr.ErrInternalServerError, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("store: RemoveUser: error: %w: %v", repoErr.ErrRecordNotFound, err)
	}
	return nil

}
