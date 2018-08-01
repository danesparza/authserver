package data_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

//	Gets the database path for this environment:
func getTestFile() string {
	return fmt.Sprintf("%s/testdatabase.db", os.TempDir())
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
	response, secret, err := db.Init()

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if response.ID != 1 || response.Name != "Admin" {
		t.Errorf("Init failed: Should create admin user: %+v", response)
	}

	if secret == "" {
		t.Errorf("Init failed: Should return admin user secret: %s", secret)
	}
}
