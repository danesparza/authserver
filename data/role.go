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
	Deleted     zero.Time   `json:"deleted"`
	DeletedBy   null.String `json:"deleted_by"`
}

// AddRole adds a role to the system
func (store DBManager) AddRole(context User, role Role) (Role, error) {
	//	Our return item
	retval := Role{}

	//	Validate:  Does the context user have permission to execute the request?
	if store.userIsSystemAdmin(context.ID) == false && store.userIsResourceDelegate(context.ID) == false {
		//	Return an error:
		return retval, fmt.Errorf("User '%s' does not have permission to add a role to the system", context.Name)
	}

	//	Start a transaction:
	tx, err := store.systemdb.Begin()
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
	err = store.systemdb.QueryRow(`select id, name, description, created, createdby, updated, updatedby, deleted, deletedby from role WHERE id=$1;`, roleID).Scan(
		&retval.ID,
		&retval.Name,
		&retval.Description,
		&retval.Created,
		&retval.CreatedBy,
		&retval.Updated,
		&retval.UpdatedBy,
		&retval.Deleted,
		&retval.DeletedBy,
	)
	if err != nil {
		return retval, fmt.Errorf("Problem fetching role: %s", err)
	}

	//	Return it:
	return retval, nil
}

// GetAllRoles returns an array of all roles
func (store DBManager) GetAllRoles(context User) ([]Role, error) {
	retval := []Role{}

	//	Get all the items:
	rows, err := store.systemdb.Query("select id, name, description, created, createdby, updated, updatedby, deleted, deletedby from role")
	if err != nil {
		return retval, fmt.Errorf("Problem selecting all roles: %s", err)
	}

	for rows.Next() {
		item := Role{}

		if err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Created,
			&item.CreatedBy,
			&item.Updated,
			&item.UpdatedBy,
			&item.Deleted,
			&item.DeletedBy); err != nil {
			rows.Close()
			break
		}

		retval = append(retval, item)
	}

	if err = rows.Err(); err != nil {
		return retval, fmt.Errorf("Problem scanning all roles: %s", err)
	}

	//	Return our slice:
	return retval, nil
}

// roleExists returns 'true' if the role can be found, 'false' if it can't be found
func (store DBManager) roleExists(role Role) bool {
	retval := false

	item := Role{}
	err := store.systemdb.QueryRow("SELECT id, name FROM role WHERE id=$1;", role.ID).Scan(
		&item.ID,
		&item.Name,
	)
	if err != nil {
		return false
	}

	//	We've gotten this far -- we must have found something
	retval = true

	return retval
}
