package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestInitialization(t *testing.T) {
	t.Logf("Using UNITTEST_INFLUX_URL: %s", os.Getenv("UNITTEST_INFLUX_URL"))
}

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

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

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

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	No items are in the database!

	//	Act
	response, err := db.GetAllUsers(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten all users without error: %s", err)
	}

	if len(response) != 0 {
		t.Errorf("GetAllUsers failed: Should not have gotten any items")
	}
}

func TestUser_GetAllUsers_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Try storing some users:
	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	})
	if err != nil {
		t.Errorf("SetUser 1 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	})
	if err != nil {
		t.Errorf("SetUser 2 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser3",
		SecretHash:  "SomeRandomSecret3",
		Description: "Unit test user 3",
	})
	if err != nil {
		t.Errorf("SetUser 3 failed: Should have created users without error: %s", err)
	}

	//	Act
	response, err := db.GetAllUsers(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllUsers failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 3 {
		t.Errorf("GetAllUsers failed: Should have gotten all users")
	}
}

func TestUser_GetUserByID_IdDoesntExist_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Try storing some users:
	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	})
	if err != nil {
		t.Errorf("SetUser 1 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	})
	if err != nil {
		t.Errorf("SetUser 2 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser3",
		SecretHash:  "SomeRandomSecret3",
		Description: "Unit test user 3",
	})
	if err != nil {
		t.Errorf("SetUser 3 failed: Should have created users without error: %s", err)
	}

	//	Act
	response, err := db.GetUserByID(uctx, 42)

	//	Assert
	if err != nil {
		t.Errorf("GetUserByID failed: Should get an item without error: %s", err)
	}

	if response.ID != 0 {
		t.Errorf("GetUserByID failed: Should fetch a blank item: %+v", response)
	}
}

func TestUser_GetUserByID_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	Try storing some users:
	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	})
	if err != nil {
		t.Errorf("SetUser 1 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	})
	if err != nil {
		t.Errorf("SetUser 2 failed: Should have created users without error: %s", err)
	}

	_, err = db.SetUser(uctx, data.User{
		Name:        "TestUser3",
		SecretHash:  "SomeRandomSecret3",
		Description: "Unit test user 3",
	})
	if err != nil {
		t.Errorf("SetUser 3 failed: Should have created users without error: %s", err)
	}

	//	Act
	response, err := db.GetUserByID(uctx, 2)

	//	Assert
	if err != nil {
		t.Errorf("GetUserByID failed: Should get an item without error: %s", err)
	}

	if response.Name != "TestUser2" {
		t.Errorf("GetUserByID failed: Should get an item with the correct name: %+v", response)
	}
}

func TestUser_GetUserByID_NoData_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Our 'context' user (the one performing the action)
	uctx := data.User{
		Name: "Admin",
	}

	//	No data exists yet!

	//	Act
	response, err := db.GetUserByID(uctx, 1)

	//	Assert
	if err != nil {
		t.Errorf("GetUserByID failed: Should get an item without error: %s", err)
	}

	if response.ID != 0 {
		t.Errorf("GetUserByID failed: Should fetch a blank item: %+v", response)
	}
}
