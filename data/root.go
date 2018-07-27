package data

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
)

// SystemDB is the BoltDB database for
// user/application/role storage
type SystemDB struct {
	db *bolt.DB
}

// TokenDB is the BoltDB database for
// token storage
type TokenDB struct {
	Database string
}

// NewSystemDB creates a new instance of a SystemDB
func NewSystemDB(filepath string) (*SystemDB, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("An error occurred opening the SystemDB: %s", err)
	}

	return &SystemDB{db}, nil
}

// Path returns the database path for the SystemDB
func (b *SystemDB) Path() string {
	return b.db.Path()
}

// Close closes the SystemDB database
func (b *SystemDB) Close() {
	b.db.Close()
}
