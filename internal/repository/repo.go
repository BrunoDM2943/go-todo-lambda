package repository

import "github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"

//go:generate mockgen -source=./repo.go -destination=./mock/repo_mock.go
type TodoRepository interface {
	Save(item *model.Item) error
	FindByID(id int64) (*model.Item, error)
	ListAll() ([]*model.Item, error)
	DeleteByID(id int64) error
}
