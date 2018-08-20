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
