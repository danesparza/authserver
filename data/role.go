package data

import (
	"time"
)

// Role defines a role or permission that a user is assigned within an
// application/role/service
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Updated     time.Time `json:"updated"`
	UpdatedBy   string    `json:"updated_by"`
	Deleted     time.Time `json:"deleted"`
	DeletedBy   string    `json:"deleted_by"`
}
