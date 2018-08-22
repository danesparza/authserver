package data

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ScopeUser is a hierarchy of a user and the resource and role
// scopes they have been assigned
type ScopeUser struct {
	ID             string
	Name           string
	Description    string
	ScopeResources []ScopeResource
}

// ScopeResource is part of the user/resource/role scope hierarchy
type ScopeResource struct {
	ID          string
	Name        string
	Description string
	ScopeRoles  []ScopeRole
}

// ScopeRole is part of the user/resource/role scope hierarchy
type ScopeRole struct {
	ID          string
	Name        string
	Description string
}

// GetUserScopesWithCredentials - verifies credentials and returns the scopeuser hierarchy
func (store DBManager) GetUserScopesWithCredentials(name, secret string) (ScopeUser, error) {
	retUser := ScopeUser{}

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

	//	If everything checks out, get the scopeuser information and return it:
	retUser, err = store.getUserScopes(user)
	if err != nil {
		return retUser, fmt.Errorf("Problem fetching scopes for the user: %s", err)
	}

	//	Return our user:
	return retUser, nil
}

// getUserScopes gets the scope hierarchy for a given user
func (store DBManager) getUserScopes(user User) (ScopeUser, error) {

	//	First, copy the necessary properties from the passed user
	retval := ScopeUser{
		ID:          user.ID,
		Name:        user.Name,
		Description: user.Description,
	}

	//	Next, look in the user_resource_role table:
	//	-- see what resources they have
	//	-- see what roles they have on those resources
	//	Build up the ScopeUser hierarchy
	rows, err := store.systemdb.Query(getResourcesForUser, user.ID)
	if err != nil {
		return retval, fmt.Errorf("Problem getting resources for user %s / %v: %s", user.Name, user.ID, err)
	}

	for rows.Next() {
		gres := ScopeResource{}

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
			grole := ScopeRole{}

			if err = rolesrows.Scan(
				&grole.ID,
				&grole.Name,
				&grole.Description); err != nil {
				rolesrows.Close()
				break
			}

			gres.ScopeRoles = append(gres.ScopeRoles, grole)

		}

		if err = rolesrows.Err(); err != nil {
			return retval, fmt.Errorf("Problem scanning roles for user %s / %v & resource %s / %v: %s", user.Name, user.ID, gres.Name, gres.ID, err)
		}

		retval.ScopeResources = append(retval.ScopeResources, gres)
	}

	if err = rows.Err(); err != nil {
		return retval, fmt.Errorf("Problem scanning resources for user %s / %v: %s", user.Name, user.ID, err)
	}

	return retval, nil
}
