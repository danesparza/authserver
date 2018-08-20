package data_test

import (
	"os"
	"testing"
	"time"

	"github.com/danesparza/authserver/data"
)

func TestToken_GetNewToken_ValidUserAndExpiresafter_Successful(t *testing.T) {
	//	Arrange
	systemdbfilename, tokendbfilename := getTestFiles()
	defer os.Remove(systemdbfilename)
	defer os.Remove(tokendbfilename)

	db, err := data.NewDBManager(systemdbfilename, tokendbfilename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	uctx, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our new password:
	userPassword := "newpassword"

	//	Add a users:
	newUser1, err := db.AddUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	}, userPassword)
	if err != nil {
		t.Errorf("AddUser failed: Should have created user1 without issue, but got error: %s", err)
	}

	//	Add user to the system / admin resource & role
	_, err = db.AddUserToResourceWithRole(uctx, newUser1, data.Resource{ID: "bdldpjad2pm0cd64ra81"}, data.Role{ID: "bdldpjad2pm0cd64ra82"})
	if err != nil {
		t.Errorf("AddUserToResourceWithRole failed: Should have added newUser1 to the system admin resource and role without error, but got: %s", err)
	}

	//	Act
	tokenResponse, err := db.GetNewToken(newUser1, 5*time.Minute)

	//	Assert
	if err != nil {
		t.Errorf("GetNewToken failed: Should have gotten token without an error, but got: %s", err)
	}

	t.Logf("Token response: %+v", tokenResponse)

}
