package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

//	Gets the database path for this environment:
func getTestFile() string {
	return "testdatabase.db"
}

func TestRoot_Init_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()
	//	No data exists yet!

	//	Act
	response, secret, err := db.AuthSystemBootstrap()

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if response.ID != "bdldpjad2pm0cd64ra80" || response.Name != "admin" {
		t.Errorf("Init failed: Should create admin user: %+v", response)
	}

	if secret == "" {
		t.Errorf("Init failed: Should return admin user secret: %s", secret)
	}

	t.Logf("New Admin user: %+v", response)
	t.Logf("New Admin user secret: %s", secret)
}
