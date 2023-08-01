package controller

import (
	"example/model"

	"github.com/stretchr/testify/mock"
)

type BookMockRepository struct {
	Mock mock.Mock
}

func (m *BookMockRepository) FindAll() ([]model.Book, error) {
	args := m.Mock.Called()
	return args[0].([]model.Book), args.Error(1)
}

func (m *BookMockRepository) FindById(id int) (model.Book, error) {
	args := m.Mock.Called(id)

	return args.Get(0).(model.Book), args.Error(1)
}

func (m *BookMockRepository) Create(b *model.Book) error {
	args := m.Mock.Called(b)
	return args.Error(0)
}

func (m *BookMockRepository) Delete(id int) error {
	args := m.Mock.Called(id)
	return args.Error(0)
}
