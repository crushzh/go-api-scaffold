package model

import "time"

// Example is the example model (CRUD demo)
// New modules can reference this model definition
type Example struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null;index"`
	Description string    `json:"description" gorm:"size:500"`
	Status      string    `json:"status" gorm:"size:20;default:active"` // active, inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName overrides the table name
func (Example) TableName() string {
	return "examples"
}

// CreateExampleRequest is the create request
type CreateExampleRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateExampleRequest is the update request
type UpdateExampleRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// QueryExampleRequest is the query request
type QueryExampleRequest struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   string `form:"status" json:"status"`
}
