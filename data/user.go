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
	if store.userHasResourceRole(context.ID, systemResourceID, systemAdminRoleID, systemDelegateRoleID) == false {
		//	Return an error:
		return retval, fmt.Errorf("User %s does not have permission to add a user to the system", context.Name)
	}

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

// GetUserGrantsWithCredentials - verifies credentials and returns the grantuser hierarchy
func (store SystemDB) GetUserGrantsWithCredentials(name, secret string) (GrantUser, error) {
	retUser := GrantUser{}

	//	First, find the user with the given name and get the hashed password
	user := User{}
	err := store.db.QueryRow("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE name=$1;", name).Scan(
		&user.ID,
		&user.Enabled,
		&user.Name,
		&user.Description,
		&user.SecretHash,
		&user.Created,
		&user.CreatedBy,
		&user.Updated,
		&user.UpdatedBy,
		&user.Deleted,
		&user.DeletedBy,
	)
	if err != nil {
		return retUser, fmt.Errorf("Problem selecting user: %s", err)
	}

	// Compare the given password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(user.SecretHash), []byte(secret))
	if err != nil { // nil means it is a match
		return retUser, fmt.Errorf("The user was not found or the password was incorrect")
	}

	//	If everything checks out, get the grantuser information and return it:
	retUser, err = store.getUserGrants(user)
	if err != nil {
		return retUser, fmt.Errorf("Problem fetching grants for the user: %s", err)
	}

	//	Return our user:
	return retUser, nil
}

// userHasResourceRole returns 'true' if a given user has a given resource role, false if they don't
func (store SystemDB) userHasResourceRole(userID, resourceID string, roleIDs ...string) bool {
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

	err := store.db.QueryRow(query, qrargs...).Scan(
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

// getUserGrants gets the grant hierarchy for a given user
func (store SystemDB) getUserGrants(user User) (GrantUser, error) {

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
	rows, err := store.db.Query(getResourcesForUser, user.ID)
	if err != nil {
		return retval, fmt.Errorf("Problem getting resources for user %s / %v: %s", user.Name, user.ID, err)
	}

	for rows.Next() {
		gres := GrantResource{}

		if err = rows.Scan(
			&gres.ID,
			&gres.Name,
			&gres.Description); err != nil {
			rows.Close()
			break
		}

		//	Now that we have a resource, see what roles we should add to it for this user:
		rolesrows, err := store.db.Query(getRolesForUserAndResources, user.ID, gres.ID)
		if err != nil {
			return retval, fmt.Errorf("Problem getting resources for user %s / %v: %s", user.Name, user.ID, err)
		}

		for rolesrows.Next() {
			grole := GrantRole{}

			if err = rolesrows.Scan(
				&grole.ID,
				&grole.Name,
				&grole.Description); err != nil {
				rolesrows.Close()
				break
			}

			gres.GrantRoles = append(gres.GrantRoles, grole)

		}

		if err = rolesrows.Err(); err != nil {
			return retval, fmt.Errorf("Problem scanning roles for user %s / %v & resource %s / %v: %s", user.Name, user.ID, gres.Name, gres.ID, err)
		}

		retval.GrantResources = append(retval.GrantResources, gres)
	}

	if err = rows.Err(); err != nil {
		return retval, fmt.Errorf("Problem scanning resources for user %s / %v: %s", user.Name, user.ID, err)
	}

	return retval, nil
}

// AddUserToResourceWithRole adds the specified user to the resource and assigns the given role.
// Returns an error if the user, resource, or role don't already exist
func (store SystemDB) AddUserToResourceWithRole(context, user User, resource Resource, role Role) (UserResourceRole, error) {

	//	Our return item
	retval := UserResourceRole{}

	//	Get the user/resource/role - make sure they all exist.
	//	Throw an error if one of them doesn't exist in the system

	//	Add / update the item in the system

	//	Return our result

	return retval, nil
}
