package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bolt "github.com/coreos/bbolt"
)

// Resource represents an application / resource / service in the system
// It is associated with users (and user roles)
type Resource struct {
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

// SetResource adds or updates a resource in the system
func (store SystemDB) SetResource(context User, resource Resource) (Resource, error) {

	//	Our return item
	retval := Resource{}

	//	Update the database:
	err := store.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("resources"))
		if err != nil {
			return fmt.Errorf("An error occurred getting the user bucket: %s", err)
		}

		// Generate ID for the resource if we're adding a new one.
		if resource.ID == 0 {
			id, err := b.NextSequence()
			if err != nil {
				return fmt.Errorf("An error occurred getting a resource id: %s", err)
			}
			resource.ID = int64(id)
		}

		//	Set the current datetime(s) and created/updated by information:
		if resource.Created.IsZero() {
			resource.Created = time.Now()
			resource.CreatedBy = context.Name
		}

		resource.Updated = time.Now()
		resource.UpdatedBy = context.Name

		//	Serialize to JSON format
		encoded, err := json.Marshal(resource)
		if err != nil {
			return err
		}

		//	Store it, with the 'id' as the key:
		keyName := strconv.FormatInt(resource.ID, 10)
		return b.Put([]byte(keyName), encoded)
	})

	//	Set our return value:
	retval = resource

	return retval, err
}

// GetAllResources returns an array of all users
func (store SystemDB) GetAllResources() ([]Resource, error) {
	retval := []Resource{}

	//	Get all the items:
	err := store.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			item := Resource{}
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("An error occurred deserializing all resources: %s", err)
			}

			//	Add to the return slice:
			retval = append(retval, item)
		}

		return nil
	})

	//	Return our slice:
	return retval, err
}
