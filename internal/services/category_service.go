package services

import "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"

type categoryStorage interface {
	GetAll() ([]model.Category, error)
}

type categoryService struct {
	categoryStorage categoryStorage
	categories      []model.Category
}

func NewCategoryService(s categoryStorage) (*categoryService, error) {
	cts, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	return &categoryService{s, cts}, nil
}

func (s *categoryService) GetAll() []model.Category {
	return s.categories
}
