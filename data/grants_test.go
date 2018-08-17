package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestGrants_GetGrantUserWithCredentials_ValidCredentials_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewDBManager(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	No data exists yet, so bootstrap
	response, secret, err := db.AuthSystemBootstrap()

	//	Act
	grantinfo, gerr := db.GetUserGrantsWithCredentials(response.Name, secret)

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if gerr != nil {
		t.Errorf("GetUserGrantsWithCredentials: Should get grants without error: %s", err)
	}

	if response.ID != "bdldpjad2pm0cd64ra80" || response.Name != "admin" {
		t.Errorf("Init failed: Should create admin user: %+v", response)
	}

	if secret == "" {
		t.Errorf("Init failed: Should return admin user secret: %s", secret)
	}

	if len(grantinfo.GrantResources) != 1 {
		t.Errorf("GetUserGrantsWithCredentials failed: Should return all grants, but got: %v", len(grantinfo.GrantResources))
	}

	//	Spit out what we found (for debugging):
	//	t.Logf("Grants found: %+v", grantinfo)
}

func TestGrants_GetGrantUserWithCredentials_WrongCredentials_ReturnsError(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewDBManager(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	No data exists yet, so bootstrap
	response, _, err := db.AuthSystemBootstrap()

	//	Act
	grantinfo, gerr := db.GetUserGrantsWithCredentials(response.Name, "INTENTIONALLY_WRONG_AND_VERY_INCORRECT_PASSWORD")

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if gerr == nil {
		t.Errorf("GetUserGrantsWithCredentials: Should return error, but didn't")
	}

	if len(grantinfo.GrantResources) != 0 {
		t.Errorf("GetUserGrantsWithCredentials failed: Should not retury any grant information, but got: %v", len(grantinfo.GrantResources))
	}

	//	Spit out what we found (for debugging):
	//	t.Logf("Grants found: %+v", grantinfo)
}
