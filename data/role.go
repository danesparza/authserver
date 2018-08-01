package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bolt "github.com/coreos/bbolt"
)

// Role defines a role or permission that a user is assigned within an
// application/role/service
type Role struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Updated     time.Time `json:"updated"`
	UpdatedBy   string    `json:"updated_by"`
	Deleted     time.Time `json:"deleted"`
	DeletedBy   string    `json:"deleted_by"`
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
	err = store.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("roles"))
		if err != nil {
			return fmt.Errorf("An error occurred getting the role bucket: %s", err)
		}

		// Generate ID for the role if we're adding a new one.
		if role.ID == 0 {
			id, err := b.NextSequence()
			if err != nil {
				return fmt.Errorf("An error occurred getting a role id: %s", err)
			}
			role.ID = int64(id)
		}

		//	Set the current datetime(s) and created/updated by information:
		if role.Created.IsZero() {
			role.Created = time.Now()
			role.CreatedBy = context.Name
		}

		role.Updated = time.Now()
		role.UpdatedBy = context.Name

		//	Serialize to JSON format
		encoded, err := json.Marshal(role)
		if err != nil {
			return err
		}

		//	Store it, with the 'id' as the key:
		keyName := strconv.FormatInt(role.ID, 10)
		return b.Put([]byte(keyName), encoded)
	})

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
	err = store.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("roles"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our item
			item := Role{}
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("An error occurred deserializing all roles: %s", err)
			}

			//	Add to the return slice:
			retval = append(retval, item)
		}

		return nil
	})

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
