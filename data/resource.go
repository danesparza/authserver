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

// SetResource adds or updates a resource in the system
func (store SystemDB) SetResource(context User, resource Resource) (Resource, error) {

	//	Our return item
	retval := Resource{}

	//	Validate:  Does the context user have permission to execute the request?

	//	Log the request:
	fields := map[string]interface{}{
		"context_name":  context.Name,
		"context_id":    context.ID,
		"resource_name": resource.Name,
		"resource_id":   resource.ID,
	}
	err := store.Log("resource_activity", "SetResource_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	If the passed resource doesn't have an id, treat it as new and add it:
	if resource.ID == "" {
		tx, err := store.db.Begin()
		if err != nil {
			return retval, fmt.Errorf("An error occurred starting a transaction for a resource: %s", err)
		}

		//	Generate a sysid:
		resourceSysID := xid.New().String()

		_, err = tx.Exec(`INSERT INTO 
			resource (sysid, name, description, created, createdby, updated, updatedby) 
			VALUES ($1, $2, $3, $4, $5, $6, $7);`,
			resourceSysID,
			resource.Name,
			resource.Description,
			resource.Created,
			resource.CreatedBy,
			resource.Updated,
			resource.UpdatedBy)
		if err != nil {
			return retval, fmt.Errorf("An error occurred adding a resource: %s", err)
		}

		err = tx.Commit()
		if err != nil {
			return retval, fmt.Errorf("An error occurred committing a transaction for a resource: %s", err)
		}
	} else {
		//	If it has an id, update it:
		tx, err := store.db.Begin()
		if err != nil {
			return retval, fmt.Errorf("An error occurred starting a transaction for a resource: %s", err)
		}

		resource.Updated = time.Now()
		resource.UpdatedBy = context.Name

		_, err = tx.Exec("UPDATE resource set name = $1 description = $2, updated = $3, updatedby = $4 where sysid = $5;", resource.Name, resource.Description, resource.Updated, resource.UpdatedBy, resource.ID)
		if err != nil {
			return retval, fmt.Errorf("An error occurred updating a resource: %s", err)
		}

		err = tx.Commit()
		if err != nil {
			return retval, fmt.Errorf("An error occurred committing a transaction for a resource: %s", err)
		}
	}

	retval = resource

	fields = map[string]interface{}{
		"context_name":  context.Name,
		"context_id":    context.ID,
		"resource_name": retval.Name,
		"resource_id":   retval.ID,
	}
	err = store.Log("resource_activity", "SetResource_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	return retval, err
}

// GetAllResources returns an array of all resources
func (store SystemDB) GetAllResources(context User) ([]Resource, error) {
	retval := []Resource{}

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
	}
	err := store.Log("resource_activity", "GetAllResources_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Get all the items:

	fields = map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"count":        len(retval),
	}
	err = store.Log("resource_activity", "GetAllResources_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Return our slice:
	return retval, err
}

//	GetResourceById - used for lookups / validation before relating data

//	GetResourceByName - used for resource checks
