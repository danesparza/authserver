package data

import (
	"fmt"
	"os"
	"time"

	// QL sql driver
	_ "github.com/cznic/ql/driver"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

// SystemDB is the BoltDB database for
// user/application/role storage
type SystemDB struct {
	db       *sqlx.DB
	ic       influxdb.Client
	hostname string
}

// TokenDB is the BoltDB database for
// token storage
type TokenDB struct {
	db *sqlx.DB
}

// NewSystemDB creates a new instance of a SystemDB
func NewSystemDB(dbpath, influxurl string) (*SystemDB, error) {
	retval := new(SystemDB)

	//	Create a reference to our bolt db
	db, err := sqlx.Connect("ql", dbpath)
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

// Close closes the SystemDB database
func (store SystemDB) Close() error {
	return store.db.Close()
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

// AuthSystemBootstrap initializes the SystemDB and creates any default admin users / roles / resources
func (store SystemDB) AuthSystemBootstrap() (User, string, error) {
	adminUser := User{}
	adminPassword := ""

	//	Start our database transaction
	tx, err := store.db.Begin()
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem starting a transaction to bootstrap auth system")
	}

	//	Create our database schema
	tx.Exec(resourceSchema)
	tx.Exec(roleSchema)
	tx.Exec(userSchema)
	tx.Exec(userIXSysID)
	tx.Exec(userIXName)

	//	Generate a password
	adminPassword = xid.New().String()

	//	Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem hashing admin password: %s", err)
	}

	//	Add our default admin user - the insert statement requires some parameters be passed:
	_, err = tx.Exec(defaultAdmin, adminID, string(hashedPassword))
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem adding admin user: %s", err)
	}

	//	Create the default system roles and resources:

	//	Commit our transaction
	err = tx.Commit()
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem committing a transaction to bootstrap auth system")
	}

	//	Get our admin user from the database and create our return object:
	adminUser = User{}
	err = store.db.Get(&adminUser, "SELECT * FROM user WHERE id=$1;", adminID)
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem fetching admin user: %s", err)
	}

	/*  For reference.  Remove if no longer needed
	systemRoles := []Role{
		{Name: "admin", Description: "Admin role:  Can create/edit/delete all users/resources/roles"},
		{Name: "editor", Description: "Editor role:  Can assign users/resources/roles"},
		{Name: "reader", Description: "Reader role:  Can view users/resources/roles"},
	}

	systemResources := []Resource{
		{Name: "authserver", Description: "Authserver resource:  Defines authserver system access"},
	}

	*/

	return adminUser, adminPassword, nil
}
