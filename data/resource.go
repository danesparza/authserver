package data

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Resource represents an application / resource / service in the system
// It is associated with users (and user roles)
type Resource struct {
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
	err = store.db.QueryRow(`SELECT id, name, description, created, createdby, updated, updatedby, deleted, deletedby FROM resource WHERE id=$1;`, resourceID).Scan(
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
		return retval, fmt.Errorf("Problem fetching resource: %s", err)
	}

	//	Return it:
	return retval, nil
}

// GetAllResources returns an array of all resources
func (store SystemDB) GetAllResources(context User) ([]Resource, error) {
	retval := []Resource{}

	//	Get all the items:
	rows, err := store.db.Query("select * from resource")
	if err != nil {
		return retval, fmt.Errorf("Problem selecting all resources: %s", err)
	}

	for rows.Next() {
		item := Resource{}

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
		return retval, fmt.Errorf("Problem scanning all resources: %s", err)
	}

	//	Return our slice:
	return retval, nil
}

//	GetResourceById - used for lookups / validation before relating data

//	GetResourceByName - used for resource checks
