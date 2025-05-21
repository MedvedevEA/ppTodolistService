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
	addMessageQuery = `
INSERT INTO message (task_id,user_id,"text") 
VALUES ($1, $2, $3) RETURNING *;`
	getMessageQuery = `
SELECT * FROM message WHERE message_id=$1;`
	getMessagesQuery = `
SELECT * FROM message 
WHERE $3::uuid IS null OR task_id = $3 
OFFSET $1 LIMIT $2;`
	updateMessageQuery = `
UPDATE message 
SET 
"text" = CASE WHEN $2::character varying IS NULL THEN "text" ELSE $2 END,
update_at = now()
WHERE message_id=$1
RETURNING *;`
	removeMessageQuery = `
DELETE FROM message 
WHERE message_id=$1;`
)

func (s *Store) AddMessage(dto *repoDto.AddMessage) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), addMessageQuery, dto.TaskId, dto.UserId, dto.Text).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.AddMessage"))
		return nil, repoErr.ErrInternalServerError
	}
	return message, nil
}
func (s *Store) GetMessage(messageId *uuid.UUID) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), getMessageQuery, messageId).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetMessage"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return message, nil
}
func (s *Store) GetMessages(dto *repoDto.GetMessages) ([]*entity.Message, error) {
	rows, err := s.pool.Query(context.Background(), getMessagesQuery, dto.Offset, dto.Limit, dto.TaskId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.GetMessages"))
		return nil, repoErr.ErrInternalServerError
	}
	defer rows.Close()
	var messages []*entity.Message
	for rows.Next() {
		message := new(entity.Message)
		err := rows.Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
		if err != nil {
			s.lg.Error(err.Error(), slog.String("owner", "store.GetMessages"))
			return nil, repoErr.ErrInternalServerError
		}
		messages = append(messages, message)
	}
	return messages, nil
}
func (s *Store) UpdateMessage(dto *repoDto.UpdateMessage) (*entity.Message, error) {
	message := new(entity.Message)
	err := s.pool.QueryRow(context.Background(), updateMessageQuery, dto.MessageId, dto.Text).Scan(&message.MessageId, &message.TaskId, &message.UserId, &message.Text, &message.CreateAt, &message.UpdateAt)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.UpdateMessage"))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoErr.ErrRecordNotFound
		}
		return nil, repoErr.ErrInternalServerError
	}
	return message, nil
}
func (s *Store) RemoveMessage(messageId *uuid.UUID) error {
	result, err := s.pool.Exec(context.Background(), removeMessageQuery, messageId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "store.RemoveMessage"))
		return repoErr.ErrInternalServerError
	}
	if result.RowsAffected() == 0 {
		s.lg.Error(sql.ErrNoRows.Error(), slog.String("owner", "store.RemoveMessage"))
		return repoErr.ErrRecordNotFound
	}
	return nil

}
