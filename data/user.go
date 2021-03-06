package data

import (
	"fmt"
	"strings"
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
	UserID     string      `json:"userid"`
	ResourceID string      `json:"resourceid"`
	RoleID     string      `json:"roleid"`
	Created    time.Time   `json:"created"`
	CreatedBy  string      `json:"created_by"`
	Updated    time.Time   `json:"updated"`
	UpdatedBy  string      `json:"updated_by"`
	Deleted    zero.Time   `json:"deleted"`
	DeletedBy  null.String `json:"deleted_by"`
}

// AddUser adds a user to the system
func (store DBManager) AddUser(context User, user User, userPassword string) (User, error) {
	//	Our return item
	retval := User{}

	//	Validate:  Does the context user have permission to make the change?
	if store.userIsSystemAdmin(context.ID) == false && store.userIsResourceDelegate(context.ID) == false {
		//	Return an error:
		return retval, fmt.Errorf("User '%s' does not have permission to add a user to the system", context.Name)
	}

	//	Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return retval, fmt.Errorf("Problem hashing user password: %s", err)
	}

	//	Generate an id:
	userID := xid.New().String()

	//	Start a transaction:
	tx, err := store.systemdb.Begin()
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
	err = store.systemdb.QueryRow("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE id=$1;", userID).Scan(
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
func (store DBManager) GetAllUsers(context User) ([]User, error) {
	retval := []User{}

	//	Get all the items:
	rows, err := store.systemdb.Query("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user")
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

// userIsSystemAdmin returns 'true' if the passed user is a system admin
func (store DBManager) userIsSystemAdmin(userID string) bool {
	retval := false

	urr := UserResourceRole{}

	//	Create the base query and suffix:
	query := "SELECT userid, resourceid, roleid, created, createdby, updated, updatedby, deleted, deletedby FROM user_resource_role WHERE userid=$1 and resourceid = $2 and roleid = $3;"
	err := store.systemdb.QueryRow(query, userID, BuiltIn.SystemResource, BuiltIn.AdminRole).Scan(
		&urr.UserID,
		&urr.ResourceID,
		&urr.RoleID,
		&urr.Created,
		&urr.CreatedBy,
		&urr.Updated,
		&urr.UpdatedBy,
		&urr.Deleted,
		&urr.DeletedBy,
	)
	if err != nil {
		return retval
	}

	//	If we got this far, we must have found the item:
	retval = true

	return retval
}

// userIsResourceDelegate returns 'true' if the passed user is a resource delegate
func (store DBManager) userIsResourceDelegate(userID string) bool {
	retval := false

	urr := UserResourceRole{}

	//	Create the base query and suffix:
	query := "SELECT userid, resourceid, roleid, created, createdby, updated, updatedby, deleted, deletedby FROM user_resource_role WHERE userid=$1 and roleid = $2;"
	err := store.systemdb.QueryRow(query, userID, BuiltIn.ResourceDelegateRole).Scan(
		&urr.UserID,
		&urr.ResourceID,
		&urr.RoleID,
		&urr.Created,
		&urr.CreatedBy,
		&urr.Updated,
		&urr.UpdatedBy,
		&urr.Deleted,
		&urr.DeletedBy,
	)
	if err != nil {
		return retval
	}

	//	If we got this far, we must have found the item:
	retval = true

	return retval
}

// userHasResourceRole returns 'true' if a given user has a given resource role, false if they don't
func (store DBManager) userHasResourceRole(userID, resourceID string, roleIDs ...string) bool {
	retval := false

	//	Sanity check -- there needs to be at least one item in roleIDs
	if len(roleIDs) < 1 {
		return false
	}

	urr := UserResourceRole{}

	//	Create our variable list of args to pass to the query
	args := []string{userID, resourceID}
	args = append(args, roleIDs...)

	//	Convert back to []interface{} (for the call to QueryRow)
	//	-- from https://stackoverflow.com/a/27689178/19020
	qrargs := make([]interface{}, len(args))
	for i, v := range args {
		qrargs[i] = v
	}

	//	Create the base query and suffix:
	query := "SELECT userid, resourceid, roleid, created, createdby, updated, updatedby, deleted, deletedby FROM user_resource_role WHERE userid=$1 and resourceid = $2 and ("
	queryroles := []string{}
	querySuffix := ");"

	//	Loop through each item in roleIDs...
	for ri := 0; ri < len(roleIDs); ri++ {
		formattedRoleIDIndex := 3 + ri
		queryroles = append(queryroles, fmt.Sprintf("roleid = $%v", formattedRoleIDIndex))
	}

	//	Append our suffix
	query = query + strings.Join(queryroles, " or ") + querySuffix

	err := store.systemdb.QueryRow(query, qrargs...).Scan(
		&urr.UserID,
		&urr.ResourceID,
		&urr.RoleID,
		&urr.Created,
		&urr.CreatedBy,
		&urr.Updated,
		&urr.UpdatedBy,
		&urr.Deleted,
		&urr.DeletedBy,
	)
	if err != nil {
		return retval
	}

	//	If we got this far, we must have found the item:
	retval = true

	return retval
}

// AddUserToResourceWithRole adds the specified user to the resource and assigns the given role.
// Returns an error if the user, resource, or role don't already exist
func (store DBManager) AddUserToResourceWithRole(context, user User, resource Resource, role Role) (UserResourceRole, error) {

	//	Our return item
	retval := UserResourceRole{}

	//	Validate:  Does the context user have permission to make the change?
	if store.userIsSystemAdmin(context.ID) == false && store.userHasResourceRole(context.ID, resource.ID, BuiltIn.ResourceDelegateRole) == false {
		//	Return an error:
		return retval, fmt.Errorf("User '%s' does not have permission to add a user to '%s/%s'", context.Name, resource.Name, role.Name)
	}

	//	Get the user/resource/role - make sure they all exist - throw an error if they don't
	if !store.userExists(user) || !store.resourceExists(resource) || !store.roleExists(role) {
		//	Throw an error if one of them doesn't exist in the system
		return retval, fmt.Errorf("The user, resource, and role must already exist in the system")
	}

	//	If they all exist, then add the item in the system...

	//	Start a transaction:
	tx, err := store.systemdb.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for a user/resource/role: %s", err)
	}

	//	Insert the item
	_, err = tx.Exec(`INSERT INTO 
		user_resource_role (userid, resourceid, roleid, created, createdby, updated, updatedby) 
			VALUES ($1, $2, $3, now(), $4, now(), $4);`,
		user.ID,
		resource.ID,
		role.ID,
		context.Name,
	)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred adding the user/resource/role: %s", err)
	}

	//	Commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for a user/resource/role: %s", err)
	}

	//	Get the user/resource/role
	err = store.systemdb.QueryRow("SELECT userid, resourceid, roleid, created, createdby, updated, updatedby, deleted, deletedby FROM user_resource_role WHERE userid=$1 and resourceid=$2 and roleid=$3;", user.ID, resource.ID, role.ID).Scan(
		&retval.UserID,
		&retval.ResourceID,
		&retval.RoleID,
		&retval.Created,
		&retval.CreatedBy,
		&retval.Updated,
		&retval.UpdatedBy,
		&retval.Deleted,
		&retval.DeletedBy,
	)
	if err != nil {
		return retval, fmt.Errorf("Problem selecting user/resource/role object: %s", err)
	}

	//	Return our result
	return retval, nil
}

// userExists returns 'true' if the user can be found, 'false' if it can't be found
func (store DBManager) userExists(user User) bool {
	retval := false

	item := User{}
	err := store.systemdb.QueryRow("SELECT id, name FROM user WHERE id=$1;", user.ID).Scan(
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

// getUserForUserID returns the user information for the given userID
func (store DBManager) getUserForUserID(userID string) (User, error) {
	item := User{}

	err := store.systemdb.QueryRow(`SELECT 
		id, enabled, name, description, created, createdby, updated, updatedby, deleted, deletedby 
		FROM user 
		WHERE id=$1;`, userID).Scan(
		&item.ID,
		&item.Enabled,
		&item.Name,
		&item.Description,
		&item.Created,
		&item.CreatedBy,
		&item.Updated,
		&item.UpdatedBy,
		&item.Deleted,
		&item.DeletedBy,
	)
	if err != nil {
		return item, fmt.Errorf("There was an error getting the user: %s", err)
	}

	return item, nil
}
