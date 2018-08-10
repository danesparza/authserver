package data

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v3"

	"gopkg.in/guregu/null.v3/zero"
)

// User represents a user in the system.  Users
// are associated with resources and roles within those applications/resources/services.
// They can be created/updated/deleted.  If they are deleted, eventually
// they will be removed from the system.  The admin user can only be disabled, not deleted
type User struct {
	ID          string      `json:"id"`
	Enabled     bool        `json:"enabled"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	SecretHash  string      `json:"secrethash"`
	Created     time.Time   `json:"created"`
	CreatedBy   string      `json:"created_by"`
	Updated     time.Time   `json:"updated"`
	UpdatedBy   string      `json:"updated_by"`
	Deleted     zero.Time   `json:"deleted"`
	DeletedBy   null.String `json:"deleted_by"`
}

// UserResourceRole defines a relationship between a user,
// a resource (application/service), and the roles that user has
// been assigned within the resource (application/service)
type UserResourceRole struct {
	UserID     string    `json:"userid"`
	ResourceID string    `json:"resourceid"`
	RoleID     string    `json:"roleid"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Updated    time.Time `json:"updated"`
	UpdatedBy  string    `json:"updated_by"`
	Deleted    time.Time `json:"deleted"`
	DeletedBy  string    `json:"deleted_by"`
}

// GrantUser is a hierarchy of a user and the resource and role
// grants they have been assigned
type GrantUser struct {
	ID             string
	Name           string
	Description    string
	GrantResources []GrantResource
}

// AddUser adds a user to the system
func (store SystemDB) AddUser(context User, user User, userPassword string) (User, error) {
	//	Our return item
	retval := User{}

	//	Validate:  Does the context user have permission to make the change?

	//	Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return retval, fmt.Errorf("Problem hashing user password: %s", err)
	}

	//	Generate an id:
	userID := xid.New().String()

	//	Start a transaction:
	tx, err := store.db.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for a user: %s", err)
	}

	//	Insert the item
	_, err = tx.Exec(`INSERT INTO 
			user (id, enabled, name, description, secrethash, created, createdby, updated, updatedby) 
			VALUES ($1, true, $2, $3, $4, now(), $5, now(), $5);`,
		userID,
		user.Name,
		user.Description,
		string(hashedPassword),
		context.Name)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred adding a user: %s", err)
	}

	//	Commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for a user: %s", err)
	}

	//	Get the user
	err = store.db.QueryRow("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE id=$1;", userID).Scan(
		&retval.ID,
		&retval.Enabled,
		&retval.Name,
		&retval.Description,
		&retval.SecretHash,
		&retval.Created,
		&retval.CreatedBy,
		&retval.Updated,
		&retval.UpdatedBy,
		&retval.Deleted,
		&retval.DeletedBy,
	)
	if err != nil {
		return retval, fmt.Errorf("Problem selecting user: %s", err)
	}

	//	Return it:
	return retval, nil
}

// GetAllUsers returns an array of all users
func (store SystemDB) GetAllUsers(context User) ([]User, error) {
	retval := []User{}

	//	Get all the items:
	rows, err := store.db.Query("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user")
	if err != nil {
		return retval, fmt.Errorf("Problem selecting all users: %s", err)
	}

	for rows.Next() {
		item := User{}

		if err = rows.Scan(
			&item.ID,
			&item.Enabled,
			&item.Name,
			&item.Description,
			&item.SecretHash,
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
		return retval, fmt.Errorf("Problem scanning all users: %s", err)
	}

	//	Return our slice:
	return retval, nil
}

// GetUserWithCredentials - used for token creation process (to login a user)
func (store SystemDB) GetUserWithCredentials(name, secret string) (User, string, error) {
	retUser := User{}
	retToken := ""

	//	Return our user:
	return retUser, retToken, nil
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

// getUserGrants gets the grant hierarchy for a given user
func (store SystemDB) getUserGrants(user User) GrantUser {

	//	First, copy the necessary properties from the passed user
	retval := GrantUser{
		ID:          user.ID,
		Name:        user.Name,
		Description: user.Description,
	}

	//	Next, look in the user_resource_role table:
	//	-- see what resources they have
	//	-- see what roles they have on those resources
	//	Build up the GrantUser hierarchy

	//

	return retval
}
