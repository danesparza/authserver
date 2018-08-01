package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	bolt "github.com/coreos/bbolt"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

// SystemDB is the BoltDB database for
// user/application/role storage
type SystemDB struct {
	db       *bolt.DB
	ic       influxdb.Client
	hostname string
}

// TokenDB is the BoltDB database for
// token storage
type TokenDB struct {
	db *bolt.DB
}

// NewSystemDB creates a new instance of a SystemDB
func NewSystemDB(dbpath, influxurl string) (*SystemDB, error) {
	retval := new(SystemDB)

	//	Create a reference to our bolt db
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("An error occurred opening the SystemDB: %s", err)
	}
	retval.db = db

	//	Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("An error occurred getting the hostname: %s", err)
	}
	retval.hostname = hostname

	//	If we have an influxurl, use it to spin up a client:
	if influxurl != "" {
		ic, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{Addr: influxurl})
		if err != nil {
			return nil, fmt.Errorf("An error occurred creating the InfluxDB client: %s", err)
		}

		//	Include the influxDB client
		retval.ic = ic
	}

	//	Return our systemdb reference
	return retval, nil
}

// Path returns the database path for the SystemDB
func (store SystemDB) Path() string {
	return store.db.Path()
}

// Close closes the SystemDB database
func (store SystemDB) Close() {
	store.db.Close()
}

// Log sends data to influx
func (store SystemDB) Log(measurement, event string, fields map[string]interface{}) error {

	//	If we appear to have an influxdb client...
	if store.ic != nil {
		bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database: "authserver",
		})
		if err != nil {
			return fmt.Errorf("An error occurred creating influx batch points: %s", err)
		}

		//	Create our tags
		tags := map[string]string{"host": store.hostname, "event": event}

		// Create a point and add to batch
		pt, err := influxdb.NewPoint(measurement, tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Problem creating a new Influx point: %s", err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := store.ic.Write(bp); err != nil {
			return fmt.Errorf("Problem writing to InfluxDB server: %s", err)
		}
	}

	return nil
}

// Init initializes the SystemDB and creates any default admin users / roles / resources
func (store SystemDB) Init() (User, string, error) {
	adminUser := User{}
	adminPassword := ""

	//	See if the admin user exists...
	store.db.View(func(tx *bolt.Tx) error {

		//	Get our bucket
		b := tx.Bucket([]byte("users"))

		if b != nil {
			//	Determine our keyname:
			userID := int64(1) //	Admin userID
			keyname := strconv.FormatInt(userID, 10)

			//	Get the data for the key:
			itemBytes := b.Get([]byte(keyname))

			if len(itemBytes) > 0 {

				//	Unmarshal data into our item
				if err := json.Unmarshal(itemBytes, &adminUser); err != nil {
					return err
				}
			}
		}

		return nil
	})

	//	The admin user doesn't exist, so create it...
	if adminUser.ID == 0 {

		//	Generate a password
		adminPassword = xid.New().String()

		//	Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			return adminUser, adminPassword, fmt.Errorf("Problem hashing password: %s", err)
		}

		//	Create the user
		adminUser = User{
			ID:         int64(1),
			Name:       "Admin",
			Enabled:    true,
			SecretHash: string(hashedPassword),
			CreatedBy:  "System",
			Created:    time.Now(),
			UpdatedBy:  "System",
			Updated:    time.Now(),
		}

		err = store.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("users"))
			if err != nil {
				return fmt.Errorf("An error occurred getting the user bucket: %s", err)
			}

			//	Serialize to JSON format
			encoded, err := json.Marshal(adminUser)
			if err != nil {
				return err
			}

			//	Store it, with the 'id' as the key:
			keyName := strconv.FormatInt(adminUser.ID, 10)
			return b.Put([]byte(keyName), encoded)
		})

		if err != nil {
			return adminUser, adminPassword, fmt.Errorf("Problem creating admin user: %s", err)
		}

	}

	//	Make sure the system roles exist (create them if they don't)
	systemRoles := []Role{
		{ID: 1, Name: "admin", Description: "Admin role:  Can create/edit/delete all users/resources/roles"},
		{ID: 2, Name: "editor", Description: "Editor role:  Can assign users/resources/roles"},
		{ID: 3, Name: "reader", Description: "Reader role:  Can view users/resources/roles"},
	}

	for r := 0; r < len(systemRoles); r++ {
		store.SetRole(adminUser, systemRoles[r])
	}

	//	Make sure the system resources exist (create them if they don't)
	systemResources := []Resource{
		{ID: 1, Name: "authserver", Description: "Authserver resource:  Defines authserver system access"},
	}

	for r := 0; r < len(systemResources); r++ {
		store.SetResource(adminUser, systemResources[r])
	}

	//	Make sure the admin / resource /roles exist (create them if they don't)

	return adminUser, adminPassword, nil
}
