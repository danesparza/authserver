package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bolt "github.com/coreos/bbolt"
)

// User represents a user in the system.  Users
// are associated with resources and roles within those applications/resources/services.
// They can be created/updated/deleted.  If they are deleted, eventually
// they will be removed from the system
type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SecretHash  string    `json:"secrethash"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Updated     time.Time `json:"updated"`
	UpdatedBy   string    `json:"updated_by"`
	Deleted     time.Time `json:"deleted"`
	DeletedBy   string    `json:"deleted_by"`
}

// UserResourceRole defines a relationship between a user,
// a resource (application/service), and the roles that user has
// been assigned within the resource (application/service)
type UserResourceRole struct {
	UserID     int64     `json:"userid"`
	ResourceID int64     `json:"resourceid"`
	RoleID     int64     `json:"roleid"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Updated    time.Time `json:"updated"`
	UpdatedBy  string    `json:"updated_by"`
	Deleted    time.Time `json:"deleted"`
	DeletedBy  string    `json:"deleted_by"`
}

// SetUser adds or updates a user in the system
func (store SystemDB) SetUser(context, user User) (User, error) {

	//	Our return item
	retval := User{}

	//	Validate:  Does the context user have permission to make the change?

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"user_name":    user.Name,
		"user_id":      user.ID,
	}
	err := store.Log("user_activity", "SetUser_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Update the database:
	err = store.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("An error occurred getting the user bucket: %s", err)
		}

		// Generate ID for the user if we're adding a new one.
		if user.ID == 0 {
			id, err := b.NextSequence()
			if err != nil {
				return fmt.Errorf("An error occurred getting a userid: %s", err)
			}
			user.ID = int64(id)
		}

		//	Set the current datetime(s) and created/updated by information:
		if user.Created.IsZero() {
			user.Created = time.Now()
			user.CreatedBy = context.Name
		}

		user.Updated = time.Now()
		user.UpdatedBy = context.Name

		//	Serialize to JSON format
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}

		//	Store it, with the 'id' as the key:
		keyName := strconv.FormatInt(user.ID, 10)
		return b.Put([]byte(keyName), encoded)
	})

	//	Set our return value:
	retval = user

	fields = map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"user_name":    retval.Name,
		"user_id":      retval.ID,
	}
	err = store.Log("user_activity", "SetUser_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	return retval, err
}

// GetAllUsers returns an array of all users
func (store SystemDB) GetAllUsers(context User) ([]User, error) {
	retval := []User{}

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
	}
	err := store.Log("user_activity", "GetAllUsers_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Get all the items:
	err = store.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("users"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			item := User{}
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("An error occurred deserializing all users: %s", err)
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
	err = store.Log("user_activity", "GetAllUsers_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Return our slice:
	return retval, err
}

//	GetUserByName - used for token creation process (to login a user)

//	GetUserById - used for lookups / validation before relating data

// AddUserToResourceRole adds the specified user to the resource role.
// Returns an error if the user, resource, or role don't already exist
func (store SystemDB) AddUserToResourceRole(context User, urr UserResourceRole) (UserResourceRole, error) {

	//	Our return item
	retval := UserResourceRole{}

	//	Get the user/resource/role - make sure they all exist.
	//	Throw an error if one of them doesn't exist in the system

	//	Create a compound key based on all 3 ids

	//	Add / update the item in the system

	//	Return our result

	return retval, nil
}
