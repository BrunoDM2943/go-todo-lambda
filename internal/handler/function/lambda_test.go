package function

import (
	"testing"

	"github.com/golang/mock/gomock"

	"context"
	"errors"
	"net/http"

	"github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"
	mock_todo "github.com/BrunoDM2943/go-todo-lambda/internal/module/todo/mock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const defaultID = "xpto"

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_todo.NewMockService(ctrl)
	handler := NewLambdaHandler(mockService)
	handler.BuildRoutes()

	t.Run("Test Get for one ID", func(t *testing.T) {

		mockService.EXPECT().GetItem(gomock.Eq(defaultID)).Return(&model.Item{}, nil)

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": defaultID,
			},
		})
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Test Get for one ID with error", func(t *testing.T) {

		mockService.EXPECT().GetItem(gomock.Eq(defaultID)).Return(&model.Item{}, errors.New("Error"))

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": defaultID,
			},
		})
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Test Get for one ID not found", func(t *testing.T) {

		mockService.EXPECT().GetItem(gomock.Eq(defaultID)).Return(nil, nil)

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": defaultID,
			},
		})
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("Test Get for one ID Bad Request", func(t *testing.T) {
		t.SkipNow()
		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": "0",
			},
		})
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Test Get for all ID - OK", func(t *testing.T) {

		mockService.EXPECT().GetItems().Return([]*model.Item{}, nil)

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api",
			PathParameters: map[string]string{
				"id": "",
			},
		})
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Test Get for all ID - Error", func(t *testing.T) {

		mockService.EXPECT().GetItems().Return(nil, errors.New("Error"))

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "GET",
			Resource:   "/todo-api",
			PathParameters: map[string]string{
				"id": "",
			},
		})
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
}

func TestDeleteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_todo.NewMockService(ctrl)
	handler := NewLambdaHandler(mockService)
	handler.BuildRoutes()

	t.Run("Test Delete ID - OK", func(t *testing.T) {

		mockService.EXPECT().DeleteItem(gomock.Eq(defaultID)).Return(nil)

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "DELETE",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": defaultID,
			},
		})
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Test Delete ID - Error", func(t *testing.T) {

		mockService.EXPECT().DeleteItem(gomock.Eq(defaultID)).Return(errors.New("Error"))

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "DELETE",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": defaultID,
			},
		})
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})

	t.Run("Test Delete ID - Bad Request", func(t *testing.T) {

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "DELETE",
			Resource:   "/todo-api/{id}",
			PathParameters: map[string]string{
				"id": "",
			},
		})
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

func TestPostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_todo.NewMockService(ctrl)
	handler := NewLambdaHandler(mockService)
	handler.BuildRoutes()

	t.Run("Test Post Item - OK", func(t *testing.T) {
		item := &model.Item{
			Title: "List",
			Text:  "Homework",
		}
		mockService.EXPECT().PostItem(gomock.Eq(item)).Return(nil)

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Resource:   "/todo-api",
			Body:       `{"title": "List", "text":"Homework"}`,
		})
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("Test Post Item - BadRequest ", func(t *testing.T) {
		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Resource:   "/todo-api",
			Body:       `{"title": "", "text":""}`,
		})
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Test Post Item - Error", func(t *testing.T) {
		item := &model.Item{
			Title: "List",
			Text:  "Homework",
		}
		mockService.EXPECT().PostItem(gomock.Eq(item)).Return(errors.New("Error"))

		response, _ := handler.HandleRequest(context.TODO(), events.APIGatewayProxyRequest{
			HTTPMethod: "POST",
			Resource:   "/todo-api",
			Body:       `{"title": "List", "text":"Homework"}`,
		})
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})

}
