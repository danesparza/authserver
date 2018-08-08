package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestRole_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := getTestFile()

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("System database file check failed: System file %s already exists, and shouldn't", filename)
	}
}

func TestRole_AddRole_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
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

	//	Create new resource:
	r1 := data.Role{
		Name:        "TestRole",
		Description: "Unit test role",
	}

	//	Act
	response, err := db.AddRole(uctx, r1)

	//	Assert
	if err != nil {
		t.Errorf("SetRole failed: Should have set an item without error: %s", err)
	}

	if response.Created.IsZero() || response.Updated.IsZero() {
		t.Errorf("SetRole failed: Should have set an item with the correct datetime: %+v", response)
	}

	if response.CreatedBy != uctx.Name {
		t.Errorf("SetRole failed: Should have set an item with the correct 'created by' user: %+v", response)
	}

	if response.UpdatedBy != uctx.Name {
		t.Errorf("SetRole failed: Should have set an item with the correct 'updated by' user: %+v", response)
	}
}

func TestRole_GetAllRoles_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
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

	//	No items are in the database!

	//	Act
	response, err := db.GetAllRoles(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllRoles failed: Should have gotten all items without error: %s", err)
	}

	if len(response) != 0 {
		t.Errorf("GetAllRoles failed: Should not have gotten any items")
	}
}

func TestRole_GetAllRoles_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewSystemDB(filename, os.Getenv("UNITTEST_INFLUX_URL"))
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

	//	Try storing some roles:
	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole1",
		Description: "Unit test role 1",
	})
	if err != nil {
		t.Errorf("SetRole 1 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole2",
		Description: "Unit test role 2",
	})
	if err != nil {
		t.Errorf("SetRole 2 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole3",
		Description: "Unit test role 3",
	})
	if err != nil {
		t.Errorf("SetRole 3 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole4",
		Description: "Unit test role 4",
	})
	if err != nil {
		t.Errorf("SetRole 4 failed: Should have created item without error: %s", err)
	}

	//	Act
	response, err := db.GetAllRoles(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllRoles failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 4 {
		t.Errorf("GetAllRoles failed: Should have gotten all items")
	}
}
