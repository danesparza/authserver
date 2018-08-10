package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestUser_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := getTestFile()

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("System database file check failed: System file %s already exists, and shouldn't", filename)
	}
}

func TestUser_AddUser_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	_, _, err = db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Create new user:
	u1 := data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret",
		Description: "Unit test user",
	}

	//	Our new password:
	userPassword := "newpassword"

	//	Act
	newUser, errAdd := db.AddUser(uctx, u1, userPassword)

	//	Assert
	if errAdd != nil {
		t.Errorf("SetUser failed: Should have added an item without error: %s", errAdd)
	}

	if newUser.Created.IsZero() || newUser.Updated.IsZero() {
		t.Errorf("SetUser failed: Should have set an item with the correct datetime: %+v", newUser)
	}

	if newUser.CreatedBy != uctx.Name {
		t.Errorf("SetUser failed: Should have set an item with the correct 'created by' user: %+v", newUser)
	}

	if newUser.UpdatedBy != uctx.Name {
		t.Errorf("SetUser failed: Should have set an item with the correct 'updated by' user: %+v", newUser)
	}
}

func TestUser_GetAllUsers_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	_, _, err = db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Only the admin is in the database!

	//	Act
	response, err := db.GetAllUsers(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten all users without error: %s", err)
	}

	if len(response) != 1 {
		t.Errorf("GetAllUsers failed: Should have gotten a single item")
	}
}

func TestUser_GetAllUsers_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	_, _, err = db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Our new password:
	userPassword := "newpassword"

	//	Try storing some users:
	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	}, userPassword)
	if err != nil {
		t.Errorf("SetUser 1 failed: Should have created users without error: %s", err)
	}

	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	}, userPassword)
	if err != nil {
		t.Errorf("SetUser 2 failed: Should have created users without error: %s", err)
	}

	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser3",
		SecretHash:  "SomeRandomSecret3",
		Description: "Unit test user 3",
	}, userPassword)
	if err != nil {
		t.Errorf("SetUser 3 failed: Should have created users without error: %s", err)
	}

	//	Act
	response, err := db.GetAllUsers(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 4 { // Don't forget bootstrapping adds the admin user
		t.Errorf("GetAllUsers failed: Should have gotten all users.  Actually got: %v", len(response))
	}
}
