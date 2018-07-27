package data

import (
	"time"
)

// Resource represents an application / resource / service in the system
// It is associated with users (and user roles)
type Resource struct {
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
