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

/*
func TestUser_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := getTestFile()

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("System database file check failed: System file %s already exists, and shouldn't", filename)
	}
}

func TestUser_Set_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db := data.SystemDB{
		Database: filename}

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Create new user:
	u1 := data.User{
		Name:        "TestUser1",
		Secret:      "SomeRandomSecret",
		Description: "Unit test user",
	}

	//	Act
	response, err := db.SetUser(uctx, u1)

	//	Assert
	if err != nil {
		t.Errorf("SetUser failed: Should have set an item without error: %s", err)
	}

	if response.Created.IsZero() || response.Updated.IsZero() {
		t.Errorf("SetUser failed: Should have set an item with the correct datetime: %+v", response)
	}

	if response.CreatedBy != uctx.Name {
		t.Errorf("SetUser failed: Should have set an item with the correct 'created by' user: %+v", response)
	}

	if response.UpdatedBy != uctx.Name {
		t.Errorf("SetUser failed: Should have set an item with the correct 'updated by' user: %+v", response)
	}
}

func TestUser_GetAllUsers_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db := data.SystemDB{
		Database: filename}

	//	No items are in the database!

	//	Act
	response, err := db.GetAllUsers()

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten all users without error: %s", err)
	}

	if len(response) != 0 {
		t.Errorf("GetAllUsers failed: Should not have gotten any items")
	}
}
*/

func TestUser_GetAllUsers_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db := data.SystemDB{
		Database: filename}

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Try storing some users:
	db.SetUser(uctx, data.User{
		Name:        "TestUser1",
		Secret:      "SomeRandomSecret1",
		Description: "Unit test user 1",
	})

	db.SetUser(uctx, data.User{
		Name:        "TestUser2",
		Secret:      "SomeRandomSecret2",
		Description: "Unit test user 2",
	})

	db.SetUser(uctx, data.User{
		Name:        "TestUser3",
		Secret:      "SomeRandomSecret3",
		Description: "Unit test user 3",
	})

	//	Act
	response, err := db.GetAllUsers()

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 3 {
		t.Errorf("GetAllUsers failed: Should have gotten all users")
	}
}
