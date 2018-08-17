package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

//	Gets the database path for this environment:
func getTestFiles() (string, string) {
	return "testsystem.db", "testtoken.db"
}

func TestRoot_AuthSystemBootstrap_Successful(t *testing.T) {
	//	Arrange
	systemdbfilename, tokendbfilename := getTestFiles()
	defer os.Remove(systemdbfilename)
	defer os.Remove(tokendbfilename)

	db, err := data.NewDBManager(systemdbfilename, tokendbfilename)
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
