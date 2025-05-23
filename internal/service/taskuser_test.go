package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"ppTodolistService/internal/entity"
	"ppTodolistService/internal/mock"
	repoDto "ppTodolistService/internal/repository/dto"
	"ppTodolistService/internal/service"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTaskUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockRepository(mockCtrl)

	lg := slog.New(slog.DiscardHandler)
	s := service.MustNew(mock, lg)

	app := fiber.New()
	app.Post("/taskusers", s.AddTaskUser)
	app.Get("/taskusers", s.GetTaskUsers)
	app.Delete("/taskusers/:TaskUserId", s.RemoveTaskUser)

	//Вспомогательные переменные
	taskUserId := uuid.New()
	userId := uuid.New()
	taskId := uuid.New()

	// POST /taskusers
	t.Run("POST /taskusers. Добавление новой записи. Успех", func(t *testing.T) {
		mock.EXPECT().
			AddTaskUser(&repoDto.AddTaskUser{TaskId: &taskId, UserId: &userId}).
			Return(&entity.TaskUser{TaskUserId: &taskUserId, TaskId: &taskId, UserId: &userId}, nil)
		body := fmt.Sprintf(`{"taskId":"%v","userId":"%v"}`, taskId, userId)

		url := "/taskusers"
		httpRequest, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 201, httpResponse.StatusCode)

		var response *entity.TaskUser
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.TaskUser{TaskUserId: &taskUserId, TaskId: &taskId, UserId: &userId}, response)
	})
	// GET /taskUsers
	t.Run("GET /TaskUsers. Получение списка записей. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetTaskUsers(&repoDto.GetTaskUsers{Offset: 0, Limit: 1, TaskId: nil, UserId: nil}).
			Return([]*entity.TaskUser{{TaskUserId: &taskUserId, TaskId: &taskId, UserId: &userId}}, nil)

		url := "/taskusers?offset=0&limit=1"
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response []*entity.TaskUser
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, []*entity.TaskUser{{TaskUserId: &taskUserId, TaskId: &taskId, UserId: &userId}}, response)
	})
	// DELETE /taskusers/:taskUserId
	t.Run("DELETE /taskuseres/:taskUserId.Удаление записи по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			RemoveTaskUser(&taskUserId).
			Return(nil)

		url := fmt.Sprintf("/taskusers/%v", taskUserId)
		httpRequest, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 204, httpResponse.StatusCode)
	})
}
