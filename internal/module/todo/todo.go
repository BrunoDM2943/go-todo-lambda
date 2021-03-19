package todo

import (
	"github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"
	"github.com/BrunoDM2943/go-todo-lambda/internal/repository"
)

//go:generate mockgen -source=./todo.go -destination=./mock/todo_mock.go
type Service interface {
	PostItem(item *model.Item) error
	GetItem(id int64) (*model.Item, error)
	GetItems() ([]*model.Item, error)
	DeleteItem(id int64) error
}

type todoService struct {
	repository repository.TodoRepository
}

func NewTodoService(repository repository.TodoRepository) Service {
	return &todoService{repository}
}

func (service *todoService) PostItem(item *model.Item) error {
	return service.repository.Save(item)
}

func (service *todoService) GetItem(id int64) (*model.Item, error) {
	return service.repository.FindByID(id)
}

func (service *todoService) GetItems() ([]*model.Item, error) {
	return service.repository.ListAll()
}

func (service *todoService) DeleteItem(id int64) error {
	return service.repository.DeleteByID(id)
}
