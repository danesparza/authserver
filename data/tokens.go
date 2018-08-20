package data

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

// Token represents an auth token
type Token struct {
	ID        string `json:"token"`
	UserID    string
	Created   time.Time
	Expires   time.Time `json:"expires"`
	Deleted   zero.Time
	DeletedBy null.String
}

// GetNewToken gets a token for the given user.  If a token already exists it expires the existing token,
// generates a new token, stores it, and returns it.  If a token doesn't already exist (or it has expired)
// it generates a new token, stores it, and returns it
func (store DBManager) GetNewToken(user User, expiresafter time.Duration) (Token, error) {

	//	Create our default return value
	retval := Token{
		ID:      xid.New().String(), // Generate a new token
		UserID:  user.ID,
		Created: time.Now(),
		Expires: time.Now().Add(expiresafter),
	}

	//	Next, set existing tokens for the user to 'expired':
	//	-- start a transaction
	tx, err := store.tokendb.Begin()
	if err != nil {
		return retval, fmt.Errorf("An error occurred starting a transaction for updating existing tokens: %s", err)
	}

	//	-- update existing items
	_, err = tx.Exec(`UPDATE tokens 
		set expires = now(), deleted = now(), deletedby = "getNewToken" 
		where userid = $1;`,
		user.ID)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred updating existing tokens: %s", err)
	}

	//	Persist the token in the database
	_, err = tx.Exec(`INSERT INTO 
		tokens(token, userid, created, expires)
		VALUES($1, $2, $3, $4);`,
		retval.ID,
		retval.UserID,
		retval.Created,
		retval.Expires)
	if err != nil {
		tx.Rollback()
		return retval, fmt.Errorf("An error occurred updating existing tokens: %s", err)
	}

	//	-- commit the transaction
	err = tx.Commit()
	if err != nil {
		return retval, fmt.Errorf("An error occurred committing a transaction for updating existing tokens: %s", err)
	}

	//	Return the token
	return retval, nil
}

// getTokenInfo returns token information for a given unexpired tokenID (or an error if it can't be found)
func (store DBManager) getTokenInfo(tokenID string) (Token, error) {

	retval := Token{}

	//	Get the token (as long as it's not expired)
	err := store.tokendb.QueryRow(`SELECT 
	token, userid, created, expires, deleted, deletedby 
	FROM tokens 
	WHERE token=$1 and expires > now();`, tokenID).Scan(
		&retval.ID,
		&retval.UserID,
		&retval.Created,
		&retval.Expires,
		&retval.Deleted,
		&retval.DeletedBy,
	)
	if err != nil {
		return retval, fmt.Errorf("Problem selecting token: %s", err)
	}

	//	Return what we found:
	return retval, nil
}

// GetGrantsForToken gets Grant information for a given token
func (store DBManager) GetGrantsForToken(tokenID string) (GrantUser, error) {

	//	Create our default return value
	retval := GrantUser{}

	//	First, get the userid for the given token
	tokenInfo, err := store.getTokenInfo(tokenID)

	if err != nil {
		return retval, fmt.Errorf("There was a problem getting token information for the token: %s", err)
	}

	//	Then get the user information for the given userID:
	userInfo, err := store.getUserForUserID(tokenInfo.UserID)

	if err != nil {
		return retval, fmt.Errorf("There was a problem getting user information for the token: %s", err)
	}

	//	Next, get the grant information for the given userid
	grantInfo, err := store.getUserGrants(userInfo)

	if err != nil {
		return retval, fmt.Errorf("There was a problem getting grant information for the token: %s", err)
	}

	retval = grantInfo

	//	Return the grant information
	return retval, nil
}
