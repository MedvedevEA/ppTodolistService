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

func TestTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockRepository(mockCtrl)

	lg := slog.New(slog.DiscardHandler)
	s := service.MustNew(mock, lg)

	app := fiber.New()
	app.Post("/tasks", s.AddTask)
	app.Get("/tasks/:taskId", s.GetTask)
	app.Get("/tasks", s.GetTasks)
	app.Patch("/tasks/:taskId", s.UpdateTask)
	app.Delete("/tasks/:taskId", s.RemoveTask)

	//Вспомогательные переменные
	taskId := uuid.New()
	statusId := uuid.New()
	title := "taskTitle"
	description := "taskDescription"

	// POST /tasks
	t.Run("POST /tasks. Добавление новой задачи. Успех", func(t *testing.T) {
		mock.EXPECT().
			AddTask(&repoDto.AddTask{StatusId: &statusId, Title: title, Description: description}).
			Return(&entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, nil)
		body := fmt.Sprintf(`{"statusId":"%v","title":"%s","description":"%s"}`, statusId, title, description)

		url := "/tasks"
		httpRequest, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 201, httpResponse.StatusCode)

		var response *entity.Task
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, response)
	})
	// GET /tasks/:taskId
	t.Run("GET /tasks/:taskId. Получение задачи по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetTask(&taskId).
			Return(&entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, nil)

		url := fmt.Sprintf("/tasks/%v", taskId)
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response *entity.Task
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, response)
	})
	// GET /tasks
	t.Run("GET /tasks. Получение списка задач. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetTasks(&repoDto.GetTasks{Offset: 0, Limit: 1, StatusId: nil}).
			Return([]*entity.Task{{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}}, nil)

		url := "/tasks?offset=0&limit=1"
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response []*entity.Task
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, []*entity.Task{{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}}, response)
	})
	// PATCH /tasks/:taskId
	t.Run("PATCH /tasks/:taskId. Изменение задачи. Успех", func(t *testing.T) {
		mock.EXPECT().
			UpdateTask(&repoDto.UpdateTask{TaskId: &taskId, StatusId: &statusId, Title: &title, Description: &description}).
			Return(&entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, nil)
		body := fmt.Sprintf(`{"statusId":"%v","title":"%s","description":"%s"}`, statusId, title, description)

		url := fmt.Sprintf("/tasks/%v", taskId)
		httpRequest, err := http.NewRequest("PATCH", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response *entity.Task
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Task{TaskId: &taskId, StatusId: &statusId, Title: title, Description: description}, response)

	})
	// DELETE /tasks/:taskId
	t.Run("DELETE /taskes/:taskId.Удаление задачи по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			RemoveTask(&taskId).
			Return(nil)

		url := fmt.Sprintf("/tasks/%v", taskId)
		httpRequest, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 204, httpResponse.StatusCode)
	})
}
