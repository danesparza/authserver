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

	//	Start a transaction:
	tx, err := store.db.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for a user: %s", err)
	}

	//	Generate an id:
	userID := xid.New().String()

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
		return retval, fmt.Errorf("An error occurred adding a user: %s", err)
	}

	//	Commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for a user: %s", err)
	}

	//	Get the user
	rows, err := store.db.Query("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE id=$1;", userID)
	if err != nil {
		return retval, fmt.Errorf("Problem selecting user: %s", err)
	}

	for rows.Next() {
		if err = rows.Scan(
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
			&retval.DeletedBy); err != nil {
			rows.Close()
			break
		}
	}

	if err = rows.Err(); err != nil {
		return retval, fmt.Errorf("Problem scanning user: %s", err)
	}

	//	Return it:
	return retval, nil
}

/*
// GetAllUsers returns an array of all users
func (store SystemDB) GetAllUsers(context User) ([]User, error) {
	retval := []User{}

	//	Get all the items:
	err := store.db.Select(&retval, "select * from user")
	if err != nil {
		return retval, fmt.Errorf("Problem fetching all users: %s", err)
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
*/
