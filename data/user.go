package data

import (
	"encoding/binary"
	"encoding/json"
	"time"

	bolt "github.com/coreos/bbolt"
)

// User represents a user in the system.  Users
// are associated with resources and roles within those applications/resources/services.
// They can be created/updated/deleted.  If they are deleted, eventually
// they will be removed from the system
type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Secret      string    `json:"secret"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Updated     time.Time `json:"updated"`
	UpdatedBy   string    `json:"updated_by"`
	Deleted     time.Time `json:"deleted"`
	DeletedBy   string    `json:"deleted_by"`
}

// UserResourceRoles defines a relationship between a user,
// a resource/application/service, and the roles that user has
// been assigned within the application/resource/service
type UserResourceRoles struct {
	UserID     int       `json:"userid"`
	ResourceID int       `json:"resourceid"`
	RoleID     int       `json:"roleid"`
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

	//	Open the database
	db, err := bolt.Open(store.Database, 600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Update the database:
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		// Generate ID for the user if needed.
		if user.ID == 0 {
			id, _ := b.NextSequence()
			user.ID = int(id)
		}

		//	Set the current datetime(s) and created/updated by information:
		if user.Created.IsZero() {
			user.Created = time.Now()
			user.CreatedBy = context.Name
		}

		user.Updated = time.Now()
		user.UpdatedBy = context.Name

		//	Serialize to JSON format
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}

		//	Store it, with the 'id' as the key:
		return b.Put(itob(user.ID), encoded)
	})

	//	Set our return value:
	retval = user

	return retval, err
}

// GetAllUsers returns an array of all users
func (store SystemDB) GetAllUsers() ([]User, error) {
	retval := []User{}

	//	Open the database:
	db, err := bolt.Open(store.Database, 0600, nil)
	defer db.Close()
	if err != nil {
		return retval, err
	}

	//	Get all the items:
	err = db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("users"))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			//	Unmarshal data into our config item
			user := User{}
			if err := json.Unmarshal(v, &user); err != nil {
				return err
			}

			//	Add to the return slice:
			retval = append(retval, user)
		}

		return nil
	})

	//	Return our slice:
	return retval, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
