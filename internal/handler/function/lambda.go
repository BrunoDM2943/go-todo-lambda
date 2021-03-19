package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"
	"github.com/BrunoDM2943/go-todo-lambda/internal/module/todo"
	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler struct {
	todoService todo.Service
}

var successResponse = events.APIGatewayProxyResponse{
	StatusCode: http.StatusOK,
}

var createdResponse = events.APIGatewayProxyResponse{
	StatusCode: http.StatusCreated,
}

func NewLambdaHandler(todoService todo.Service) *lambdaHandler {
	return &lambdaHandler{
		todoService: todoService,
	}
}

func (handler *lambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	response := buildErrorResponse("Not implemented", http.StatusNotImplemented)
	switch request.HTTPMethod {
	case "GET":
		response = handler.getHandler(request)
	case "POST":
		response = handler.postHandler(request)
	case "DELETE":
		response = handler.deleteHandler(request)
	}
	return response
}

func (handler *lambdaHandler) deleteHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	id := request.PathParameters["id"]
	if id == "" {
		return buildErrorResponse("Invalid ID", http.StatusBadRequest)
	}
	if err := handler.todoService.DeleteItem(id); err != nil {
		return buildErrorResponse(err.Error(), http.StatusInternalServerError)
	}
	return successResponse
}

func (handler *lambdaHandler) postHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	item := &model.Item{}
	_ = json.Unmarshal([]byte(request.Body), item)

	if item.Title == "" || item.Text == "" {
		return buildErrorResponse("Invalid body", http.StatusBadRequest)
	}
	if err := handler.todoService.PostItem(item); err != nil {
		return buildErrorResponse(err.Error(), http.StatusInternalServerError)
	}
	return createdResponse
}

func (handler *lambdaHandler) getHandler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	id := request.PathParameters["id"]
	if id == "" {
		items, err := handler.todoService.GetItems()
		if err != nil {
			return buildErrorResponse(err.Error(), http.StatusInternalServerError)
		}
		body, _ := json.Marshal(items)
		return buildSuccessResponse(string(body))
	} else {
		item, err := handler.todoService.GetItem(id)
		if item == nil {
			return buildErrorResponse(fmt.Sprintf("ID %s not found", id), http.StatusNotFound)
		} else if err != nil {
			return buildErrorResponse(err.Error(), http.StatusInternalServerError)
		}
		body, _ := json.Marshal(item)
		return buildSuccessResponse(string(body))
	}
}

func buildErrorResponse(message string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(`{"message":"%s"}`, message),
		StatusCode: statusCode,
	}
}

func buildSuccessResponse(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: http.StatusOK,
	}
}
