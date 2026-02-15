package store

import (
	"go-api-scaffold/internal/model"

	"gorm.io/gorm"
)

// ExampleRepository is the example data repository
type ExampleRepository struct {
	db *gorm.DB
}

func NewExampleRepository(s *Store) *ExampleRepository {
	return &ExampleRepository{db: s.DB()}
}

// Create creates an example
func (r *ExampleRepository) Create(item *model.Example) error {
	return r.db.Create(item).Error
}

// FindByID returns an example by ID
func (r *ExampleRepository) FindByID(id uint) (*model.Example, error) {
	var item model.Example
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// List returns a paginated list of examples
func (r *ExampleRepository) List(page, pageSize int, keyword, status string) ([]model.Example, int64, error) {
	var items []model.Example
	var total int64

	query := r.db.Model(&model.Example{})

	// Filter conditions
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Update updates an example
func (r *ExampleRepository) Update(item *model.Example) error {
	return r.db.Save(item).Error
}

// Delete removes an example by ID
func (r *ExampleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Example{}, id).Error
}
