package data

import (
	"fmt"
	"os"
	"time"

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
	//	Create a reference to our bolt db
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("An error occurred opening the SystemDB: %s", err)
	}

	//	Get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("An error occurred getting the hostname: %s", err)
	}

	//	If we have an influxurl, use it to spin up a client:
	if influxurl != "" {
		ic, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{Addr: influxurl})
		if err != nil {
			return nil, fmt.Errorf("An error occurred creating the InfluxDB client: %s", err)
		}

		//	Include the influxDB client
		return &SystemDB{
			db:       db,
			ic:       ic,
			hostname: hostname,
		}, nil
	}

	//	Don't include the influxdb client by default
	return &SystemDB{
		db:       db,
		hostname: hostname,
	}, nil
}

// Path returns the database path for the SystemDB
func (b *SystemDB) Path() string {
	return b.db.Path()
}

// Close closes the SystemDB database
func (b *SystemDB) Close() {
	b.db.Close()
}

// Log sends data to influx
func (b *SystemDB) Log(measurement, event string, fields map[string]interface{}) error {

	//	If we appear to have an influxdb client...
	if b.ic != nil {
		bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database: "authserver",
		})
		if err != nil {
			return fmt.Errorf("An error occurred creating influx batch points: %s", err)
		}

		//	Create our tags
		tags := map[string]string{"host": b.hostname, "event": event}

		// Create a point and add to batch
		pt, err := influxdb.NewPoint(measurement, tags, fields, time.Now())
		if err != nil {
			return fmt.Errorf("Problem creating a new Influx point: %s", err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := b.ic.Write(bp); err != nil {
			return fmt.Errorf("Problem writing to InfluxDB server: %s", err)
		}
	}

	return nil
}
