package todo

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/BrunoDM2943/go-todo-lambda/internal/constants/model"
	mock_repository "github.com/BrunoDM2943/go-todo-lambda/internal/repository/mock"
)

var item = &model.Item{}
var allItems = []*model.Item{
	{}, {},
}

const defaultID = "XPTO"

func TestPostItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	service := NewTodoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().Save(item).Return(nil)
		err := service.PostItem(item)
		assert.Nil(t, err)
	})

	t.Run("Fail", func(t *testing.T) {
		mockRepo.EXPECT().Save(item).Return(errors.New("Error"))
		err := service.PostItem(item)
		assert.NotNil(t, err)
	})
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	service := NewTodoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(defaultID).Return(item, nil)
		foundItem, err := service.GetItem(defaultID)
		assert.Nil(t, err)
		assert.Equal(t, item, foundItem)
	})

	t.Run("Fail - Not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(defaultID).Return(nil, nil)
		foundItem, err := service.GetItem(defaultID)
		assert.Nil(t, err)
		assert.Nil(t, foundItem)
	})

	t.Run("Fail - Error", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(defaultID).Return(nil, errors.New("Error"))
		foundItem, err := service.GetItem(defaultID)
		assert.NotNil(t, err)
		assert.Nil(t, foundItem)
	})
}

func TestDeleteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	service := NewTodoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().DeleteByID(defaultID).Return(nil)
		assert.Nil(t, service.DeleteItem(defaultID))
	})

	t.Run("Fail", func(t *testing.T) {
		mockRepo.EXPECT().DeleteByID(defaultID).Return(errors.New("Error"))
		assert.NotNil(t, service.DeleteItem(defaultID))
	})
}

func TestGetItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockTodoRepository(ctrl)
	service := NewTodoService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().ListAll().Return(allItems, nil)
		items, err := service.GetItems()
		assert.Nil(t, err)
		assert.Equal(t, allItems, items)
	})

	t.Run("Fail", func(t *testing.T) {
		mockRepo.EXPECT().ListAll().Return(nil, errors.New("Error"))
		items, err := service.GetItems()
		assert.NotNil(t, err)
		assert.Nil(t, items)
	})
}
