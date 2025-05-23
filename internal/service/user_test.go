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

func TestUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mock.NewMockRepository(mockCtrl)

	lg := slog.New(slog.DiscardHandler)
	s := service.MustNew(mock, lg)

	app := fiber.New()

	app.Post("/users", s.AddUser)
	app.Get("/users", s.GetUsers)
	app.Delete("/users/:userId", s.RemoveUser)

	//Вспомогательные переменные
	userId := uuid.New()
	userName := "userName"
	// POST /users
	t.Run("POST /users. Добавление новой записи. Успех", func(t *testing.T) {
		mock.EXPECT().
			AddUserWithUserId(&repoDto.AddUser{UserId: &userId, Name: userName}).
			Return(&entity.User{UserId: &userId, Name: userName}, nil)
		body := fmt.Sprintf(`{"userId":"%v","name":"%v"}`, userId, userName)
		url := "/users"
		httpRequest, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
		assert.NoError(t, err)
		httpRequest.Header.Set("Content-Type", "application/json")

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 201, httpResponse.StatusCode)

		var response *entity.User
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, &entity.User{UserId: &userId, Name: userName}, response)
	})
	// GET /users
	t.Run("GET /users. Получение списка пользователей. Успех", func(t *testing.T) {
		mock.EXPECT().
			GetUsers(&repoDto.GetUsers{Offset: 0, Limit: 1, Name: &userName}).
			Return([]*entity.User{{UserId: &userId, Name: userName}}, nil)

		url := fmt.Sprintf("/users?offset=0&limit=1&name=%s", userName)
		httpRequest, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 200, httpResponse.StatusCode)

		var response []*entity.User
		json.NewDecoder(httpResponse.Body).Decode(&response)
		assert.Equal(t, []*entity.User{{UserId: &userId, Name: userName}}, response)
	})
	// DELETE /users/:userId
	t.Run("DELETE /users/:userId.Удаление записи по идентификатору. Успех", func(t *testing.T) {
		mock.EXPECT().
			RemoveUser(&userId).
			Return(nil)

		url := fmt.Sprintf("/users/%v", userId)
		httpRequest, err := http.NewRequest("DELETE", url, nil)
		assert.NoError(t, err)

		httpResponse, err := app.Test(httpRequest)
		assert.NoError(t, err)
		assert.Equal(t, 204, httpResponse.StatusCode)
	})
}
