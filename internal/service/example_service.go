package service

import (
	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/store"
)

// ExampleService handles example business logic
type ExampleService struct {
	repo *store.ExampleRepository
}

func NewExampleService(db *store.Store) *ExampleService {
	return &ExampleService{
		repo: store.NewExampleRepository(db),
	}
}

// Create creates an example
func (s *ExampleService) Create(req *model.CreateExampleRequest) (*model.Example, error) {
	item := &model.Example{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if item.Status == "" {
		item.Status = "active"
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

// GetByID returns an example by ID
func (s *ExampleService) GetByID(id uint) (*model.Example, error) {
	return s.repo.FindByID(id)
}

// List returns a paginated list of examples
func (s *ExampleService) List(req *model.QueryExampleRequest) ([]model.Example, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	return s.repo.List(req.Page, req.PageSize, req.Keyword, req.Status)
}

// Update updates an example
func (s *ExampleService) Update(id uint, req *model.UpdateExampleRequest) (*model.Example, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Status != nil {
		item.Status = *req.Status
	}

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

// Delete removes an example
func (s *ExampleService) Delete(id uint) error {
	return s.repo.Delete(id)
}
