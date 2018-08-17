package data

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GrantUser is a hierarchy of a user and the resource and role
// grants they have been assigned
type GrantUser struct {
	ID             string
	Name           string
	Description    string
	GrantResources []GrantResource
}

// GrantResource is part of the user/resource/role grant hierarchy
type GrantResource struct {
	ID          string
	Name        string
	Description string
	GrantRoles  []GrantRole
}

// GrantRole is part of the user/resource/role grant hierarchy
type GrantRole struct {
	ID          string
	Name        string
	Description string
}

// GetUserGrantsWithCredentials - verifies credentials and returns the grantuser hierarchy
func (store DBManager) GetUserGrantsWithCredentials(name, secret string) (GrantUser, error) {
	retUser := GrantUser{}

	//	First, find the user with the given name and get the hashed password
	user := User{}
	err := store.systemdb.QueryRow("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE name=$1;", name).Scan(
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

// getUserGrants gets the grant hierarchy for a given user
func (store DBManager) getUserGrants(user User) (GrantUser, error) {

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
	rows, err := store.systemdb.Query(getResourcesForUser, user.ID)
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
		rolesrows, err := store.systemdb.Query(getRolesForUserAndResources, user.ID, gres.ID)
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
