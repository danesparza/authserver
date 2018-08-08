package data

import (
	"fmt"
	"time"

	_ "github.com/cznic/ql/driver"
	"github.com/rs/xid"
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Resource represents an application / resource / service in the system
// It is associated with users (and user roles)
type Resource struct {
	ID          string      `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	Created     time.Time   `db:"created" json:"created"`
	CreatedBy   string      `db:"createdby" json:"created_by"`
	Updated     time.Time   `db:"updated" json:"updated"`
	UpdatedBy   string      `db:"updatedby" json:"updated_by"`
	Deleted     zero.Time   `db:"deleted" json:"deleted"`
	DeletedBy   null.String `db:"deletedby" json:"deleted_by"`
}

// AddResource adds a resource to the system
func (store SystemDB) AddResource(context User, resource Resource) (Resource, error) {
	//	Our return item
	retval := Resource{}

	//	Validate:  Does the context user have permission to execute the request?

	//	Start a transaction:
	tx, err := store.db.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for a resource: %s", err)
	}

	//	Generate a sysid:
	resourceID := xid.New().String()

	//	Insert the item
	_, err = tx.Exec(`INSERT INTO 
			resource (id, name, description, created, createdby, updated, updatedby) 
			VALUES ($1, $2, $3, now(), $4, now(), $4);`,
		resourceID,
		resource.Name,
		resource.Description,
		context.Name)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred adding a resource: %s", err)
	}

	//	Commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for a resource: %s", err)
	}

	//	Get the resource
	err = store.db.Get(&retval, "SELECT * FROM resource WHERE id=$1;", resourceID)
	if err != nil {
		return retval, fmt.Errorf("Problem fetching resource: %s", err)
	}

	//	Return it:
	return retval, nil
}

// GetAllResources returns an array of all resources
func (store SystemDB) GetAllResources(context User) ([]Resource, error) {
	retval := []Resource{}

	//	Get all the items:
	err := store.db.Select(&retval, "select * from resource")
	if err != nil {
		return retval, fmt.Errorf("Problem fetching all resources: %s", err)
	}

	//	Return our slice:
	return retval, nil
}

//	GetResourceById - used for lookups / validation before relating data

//	GetResourceByName - used for resource checks
