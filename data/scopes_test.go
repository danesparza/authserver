package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestScopes_GetUserScopesWithCredentials_ValidCredentials_Successful(t *testing.T) {
	//	Arrange
	systemdbfilename, tokendbfilename := getTestFiles()
	defer os.Remove(systemdbfilename)
	defer os.Remove(tokendbfilename)

	db, err := data.NewDBManager(systemdbfilename, tokendbfilename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	No data exists yet, so bootstrap
	response, secret, err := db.AuthSystemBootstrap()

	//	Act
	scopeinfo, gerr := db.GetUserScopesWithCredentials(response.Name, secret)

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if gerr != nil {
		t.Errorf("GetUserScopesWithCredentials: Should get scopes without error: %s", err)
	}

	if response.ID != "bdldpjad2pm0cd64ra80" || response.Name != "admin" {
		t.Errorf("Init failed: Should create admin user: %+v", response)
	}

	if secret == "" {
		t.Errorf("Init failed: Should return admin user secret: %s", secret)
	}

	if len(scopeinfo.ScopeResources) != 1 {
		t.Errorf("GetUserScopesWithCredentials failed: Should return all scopes, but got: %v", len(scopeinfo.ScopeResources))
	}

	//	Spit out what we found (for debugging):
	//	t.Logf("Scopes found: %+v", scopeinfo)
}

func TestScopes_GetUserScopesWithCredentials_WrongCredentials_ReturnsError(t *testing.T) {
	//	Arrange
	systemdbfilename, tokendbfilename := getTestFiles()
	defer os.Remove(systemdbfilename)
	defer os.Remove(tokendbfilename)

	db, err := data.NewDBManager(systemdbfilename, tokendbfilename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	No data exists yet, so bootstrap
	response, _, err := db.AuthSystemBootstrap()

	//	Act
	scopeinfo, gerr := db.GetUserScopesWithCredentials(response.Name, "INTENTIONALLY_WRONG_AND_VERY_INCORRECT_PASSWORD")

	//	Assert
	if err != nil {
		t.Errorf("Init failed: Should init without error: %s", err)
	}

	if gerr == nil {
		t.Errorf("GetUserScopesWithCredentials: Should return error, but didn't")
	}

	if len(scopeinfo.ScopeResources) != 0 {
		t.Errorf("GetUserScopesWithCredentials failed: Should not retury any scope information, but got: %v", len(scopeinfo.ScopeResources))
	}

	//	Spit out what we found (for debugging):
	//	t.Logf("Scopes found: %+v", scopeinfo)
}
