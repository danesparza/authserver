package data

import (
	"fmt"
	"time"

	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Role defines a role or permission that a user is assigned within an
// application/role/service
type Role struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Created     time.Time   `json:"created"`
	CreatedBy   string      `json:"created_by"`
	Updated     time.Time   `json:"updated"`
	UpdatedBy   string      `json:"updated_by"`
	Deleted     zero.Time   `db:"deleted" json:"deleted"`
	DeletedBy   null.String `db:"deletedby" json:"deleted_by"`
}

// SetRole adds or updates a role in the system
func (store SystemDB) SetRole(context User, role Role) (Role, error) {

	//	Our return item
	retval := Role{}

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"role_name":    role.Name,
		"role_id":      role.ID,
	}
	err := store.Log("role_activity", "SetRole_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Update the database:

	//	Set our return value:
	retval = role

	fields = map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"role_name":    retval.Name,
		"role_id":      retval.ID,
	}
	err = store.Log("role_activity", "SetRole_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	return retval, err
}

// GetAllRoles returns an array of all roles
func (store SystemDB) GetAllRoles(context User) ([]Role, error) {
	retval := []Role{}

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
	}
	err := store.Log("role_activity", "GetAllRoles_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Get all the items:

	fields = map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"count":        len(retval),
	}
	err = store.Log("role_activity", "GetAllRoles_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Return our slice:
	return retval, err
}

//	GetRoleById - used for lookups / validation before relating data

//	GetRoleByName - used for role checks
