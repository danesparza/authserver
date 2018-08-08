package data

import (
	"fmt"
	"time"

	"github.com/rs/xid"
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

// AddRole adds a role to the system
func (store SystemDB) AddRole(context User, role Role) (Role, error) {
	//	Our return item
	retval := Role{}

	//	Validate:  Does the context user have permission to execute the request?

	//	Start a transaction:
	tx, err := store.db.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for a role: %s", err)
	}

	//	Generate an id:
	roleID := xid.New().String()

	//	Insert the item
	_, err = tx.Exec(`INSERT INTO 
			role (id, name, description, created, createdby, updated, updatedby) 
			VALUES ($1, $2, $3, now(), $4, now(), $4);`,
		roleID,
		role.Name,
		role.Description,
		context.Name)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred adding a role: %s", err)
	}

	//	Commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for a role: %s", err)
	}

	//	Get the role
	err = store.db.Get(&retval, "SELECT * FROM role WHERE id=$1;", roleID)
	if err != nil {
		return retval, fmt.Errorf("Problem fetching role: %s", err)
	}

	//	Return it:
	return retval, nil
}

// GetAllRoles returns an array of all roles
func (store SystemDB) GetAllRoles(context User) ([]Role, error) {
	retval := []Role{}

	//	Get all the items:
	err := store.db.Select(&retval, "select * from role")
	if err != nil {
		return retval, fmt.Errorf("Problem fetching all roles: %s", err)
	}

	//	Return our slice:
	return retval, nil
}

//	GetRoleById - used for lookups / validation before relating data

//	GetRoleByName - used for role checks
