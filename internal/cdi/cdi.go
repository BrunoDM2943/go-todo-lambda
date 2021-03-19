package cdi

import (
	"github.com/BrunoDM2943/go-todo-lambda/internal/module/todo"
	"github.com/BrunoDM2943/go-todo-lambda/internal/repository"
)

var todoService todo.Service

func GetTodoService() todo.Service {
	if todoService == nil {
		todoService = todo.NewTodoService(repository.NewDynamoDB())
	}
	return todoService
}
