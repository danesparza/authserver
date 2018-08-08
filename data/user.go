package data

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"

	"gopkg.in/guregu/null.v3/zero"
)

// User represents a user in the system.  Users
// are associated with resources and roles within those applications/resources/services.
// They can be created/updated/deleted.  If they are deleted, eventually
// they will be removed from the system.  The admin user can only be disabled, not deleted
type User struct {
	ID          string      `db:"id" json:"id"`
	Enabled     bool        `db:"enabled" json:"enabled"`
	Name        string      `db:"name" json:"name"`
	Description string      `db:"description" json:"description"`
	SecretHash  string      `db:"secrethash" json:"secrethash"`
	Created     time.Time   `db:"created" json:"created"`
	CreatedBy   string      `db:"createdby" json:"created_by"`
	Updated     time.Time   `db:"updated" json:"updated"`
	UpdatedBy   string      `db:"updatedby" json:"updated_by"`
	Deleted     zero.Time   `db:"deleted" json:"deleted"`
	DeletedBy   null.String `db:"deletedby" json:"deleted_by"`
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

	//	If the passed user doesn't have an id, treat it as new and add it:
	if user.ID == "" {
		tx, err := store.db.Begin()
		if err != nil {
			return retval, fmt.Errorf("An error occurred starting a transaction for a user: %s", err)
		}

		user.Created = time.Now()
		user.CreatedBy = context.Name
		user.Updated = time.Now()
		user.UpdatedBy = context.Name

		_, err = tx.Exec("INSERT INTO user (enabled, name, description, secrethash, created, createdby, updated, updatedby) VALUES (true, $1, $2, $3, $4, $5, $6, $7);", user.Name, user.Description, user.SecretHash, user.Created, user.CreatedBy, user.Updated, user.UpdatedBy)
		if err != nil {
			return retval, fmt.Errorf("An error occurred adding a user: %s", err)
		}

		err = tx.Commit()
		if err != nil {
			return retval, fmt.Errorf("An error occurred committing a transaction for a user: %s", err)
		}
	} else {
		//	If it has an id, update it:
		tx, err := store.db.Begin()
		if err != nil {
			return retval, fmt.Errorf("An error occurred starting a transaction for a user: %s", err)
		}

		user.Updated = time.Now()
		user.UpdatedBy = context.Name

		_, err = tx.Exec("UPDATE user set name = $1 description = $2, secrethash = $3, updated = $4, updatedby = $5 where sysid = $6;", user.Name, user.Description, user.SecretHash, user.Updated, user.UpdatedBy, user.ID)
		if err != nil {
			return retval, fmt.Errorf("An error occurred updating a user: %s", err)
		}

		err = tx.Commit()
		if err != nil {
			return retval, fmt.Errorf("An error occurred committing a transaction for a user: %s", err)
		}
	}

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

	if err != nil {
		return retval, fmt.Errorf("An error occurred getting all users: %s", err)
	}

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

// GetUserWithCredentials - used for token creation process (to login a user)
func (store SystemDB) GetUserWithCredentials(name, secret string) (User, string, error) {
	retUser := User{}
	retToken := ""

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": "system",
		"context_id":   0,
		"user_name":    name,
	}
	err := store.Log("user_activity", "GetUserWithCredentials_Request", fields)
	if err != nil {
		return retUser, retToken, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Get all the items:

	if err != nil {
		return retUser, retToken, fmt.Errorf("An error occurred getting user: %s", err)
	}

	fields = map[string]interface{}{
		"context_name": "system",
		"context_id":   0,
		"user_name":    retUser.Name,
		"user_id":      retUser.ID,
	}
	err = store.Log("user_activity", "GetUserWithCredentials_Response", fields)
	if err != nil {
		return retUser, retToken, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Return our user:
	return retUser, retToken, err
}

// GetUserByID - used for lookups / validation before relating data
func (store SystemDB) GetUserByID(context User, userID int64) (User, error) {
	retval := User{}

	//	Log the request:
	fields := map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"user_id":      userID,
	}
	err := store.Log("user_activity", "GetUserByID_Request", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Open a read-only view to the data

	fields = map[string]interface{}{
		"context_name": context.Name,
		"context_id":   context.ID,
		"user_id":      retval.ID,
		"user_name":    retval.Name,
	}
	err = store.Log("user_activity", "GetUserByID_Response", fields)
	if err != nil {
		return retval, fmt.Errorf("An error occurred logging data: %s", err)
	}

	//	Return our item:
	return retval, nil
}

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
