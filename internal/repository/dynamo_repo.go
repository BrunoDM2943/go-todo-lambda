package repository

import "github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"

type dynamoDBRepo struct {
}

func NewDynamoDB() TodoRepository {
	return &dynamoDBRepo{}
}

func (repo *dynamoDBRepo) Save(item *model.Item) error {
	return nil
}
func (repo *dynamoDBRepo) FindByID(id int64) (*model.Item, error) {
	return nil, nil
}
func (repo *dynamoDBRepo) ListAll() ([]*model.Item, error) {
	return nil, nil
}
func (repo *dynamoDBRepo) DeleteByID(id int64) error {
	return nil
}
