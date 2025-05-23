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

func TestStatus(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockRepository(mockCtrl)

	lg := slog.New(slog.DiscardHandler)
	s := service.MustNew(mock, lg)

	app := fiber.New()
	app.Post("/statuses", s.AddStatus)
	app.Get("/statuses/:statusId", s.GetStatus)
	app.Get("/statuses", s.GetStatuses)
	app.Patch("/statuses/:statusId", s.UpdateStatus)
	app.Delete("/statuses/:statusId", s.RemoveStatus)

	//Вспомогательные переменные
	statusId := uuid.New()
	statusName := "statusName"

	// POST /statuses
	t.Run("POST /statuses. Добавление нового статуса. Успех", func(t *testing.T) {
		mock.EXPECT().
			AddStatus(statusName).
			Return(&entity.Status{StatusId: &statusId, Name: statusName}, nil)
		body := fmt.Sprintf(`{"name":"%s"}`, statusName)

		url := "/statuses"
		httpRequest, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 201, httpResponse.StatusCode)

		var response *entity.Status
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Status{StatusId: &statusId, Name: statusName}, response)

	})
	// GET /statuses/:statusId
	t.Run("GET /statuses/:statusId. Получение статуса по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetStatus(&statusId).
			Return(&entity.Status{StatusId: &statusId, Name: statusName}, nil)

		url := fmt.Sprintf("/statuses/%v", statusId)
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response *entity.Status
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Status{StatusId: &statusId, Name: statusName}, response)
	})
	// GET /statuses
	t.Run("GET /statuses. Получение списка статусов. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetStatuses().
			Return([]*entity.Status{{StatusId: &statusId, Name: statusName}}, nil)

		url := "/statuses"
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response []*entity.Status
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, []*entity.Status{{StatusId: &statusId, Name: statusName}}, response)
	})
	// PATCH /statuses/:statusId
	t.Run("PATCH /statuses/:statusId. Изменение статуса. Успех", func(t *testing.T) {
		mock.EXPECT().
			UpdateStatus(&repoDto.UpdateStatus{StatusId: &statusId, Name: &statusName}).
			Return(&entity.Status{StatusId: &statusId, Name: statusName}, nil)
		body := fmt.Sprintf(`{"name":"%s"}`, statusName)

		url := fmt.Sprintf("/statuses/%v", statusId)
		httpRequest, err := http.NewRequest("PATCH", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response *entity.Status
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.Status{StatusId: &statusId, Name: statusName}, response)
	})
	// DELETE /statuses/:statusId
	t.Run("DELETE /statuses/:statusId.Удаление статуса по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			RemoveStatus(&statusId).
			Return(nil)

		url := fmt.Sprintf("/statuses/%v", statusId)
		httpRequest, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 204, httpResponse.StatusCode)
	})
}
