package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestRole_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	systemdbfilename, tokendbfilename := getTestFiles()

	//	Act

	//	Assert
	if _, err := os.Stat(systemdbfilename); err == nil {
		t.Errorf("System database file check failed: System db file %s already exists, and shouldn't", systemdbfilename)
	}

	if _, err := os.Stat(tokendbfilename); err == nil {
		t.Errorf("Token database file check failed: Token db file %s already exists, and shouldn't", tokendbfilename)
	}
}

func TestRole_AddRole_Successful(t *testing.T) {
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
	bresponse, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := bresponse

	//	Create new resource:
	r1 := data.Role{
		Name:        "TestRole",
		Description: "Unit test role",
	}

	//	Act
	response, err := db.AddRole(uctx, r1)

	//	Assert
	if err != nil {
		t.Errorf("AddRole failed: Should have added an item without error: %s", err)
	}

	if response.Created.IsZero() || response.Updated.IsZero() {
		t.Errorf("AddRole failed: Should have added an item with the correct datetime: %+v", response)
	}

	if response.CreatedBy != uctx.Name {
		t.Errorf("AddRole failed: Should have added an item with the correct 'created by' user: %+v", response)
	}

	if response.UpdatedBy != uctx.Name {
		t.Errorf("AddRole failed: Should have added an item with the correct 'updated by' user: %+v", response)
	}
}

func TestRole_AddDuplicateRole_Successful(t *testing.T) {
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
	bresponse, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := bresponse

	//	Create new resource:
	r1 := data.Role{
		Name:        "TestRole",
		Description: "Unit test role",
	}

	r2 := data.Role{
		Name:        "TestRole",
		Description: "Unit test role (duplicate)",
	}

	//	Act
	_, err1 := db.AddRole(uctx, r1)
	_, err2 := db.AddRole(uctx, r2)

	//	Assert
	if err1 != nil {
		t.Errorf("AddRole failed: Should have added an item without error: %s", err)
	}

	if err2 != nil {
		t.Errorf("AddRole failed: Should have added a duplicate item without error: %s", err)
	}

}

func TestRole_GetAllRoles_NoItems_NoErrors(t *testing.T) {
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

	//	No items are in the database!

	//	Act
	response, err := db.GetAllRoles(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllRoles failed: Should have gotten all items without error: %s", err)
	}

	if len(response) != 2 { // Bootstrap adds two system roles
		t.Errorf("GetAllRoles failed: Should have fetched 2 items but got: %v", len(response))
	}
}

func TestRole_GetAllRoles_ItemsInDB_ReturnsItems(t *testing.T) {
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

	//	Try storing some roles:
	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole1",
		Description: "Unit test role 1",
	})
	if err != nil {
		t.Errorf("AddRole 1 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole2",
		Description: "Unit test role 2",
	})
	if err != nil {
		t.Errorf("AddRole 2 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole3",
		Description: "Unit test role 3",
	})
	if err != nil {
		t.Errorf("AddRole 3 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddRole(uctx, data.Role{
		Name:        "TestRole4",
		Description: "Unit test role 4",
	})
	if err != nil {
		t.Errorf("AddRole 4 failed: Should have created item without error: %s", err)
	}

	//	Act
	response, err := db.GetAllRoles(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllRoles failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 6 { // Bootstrap adds 2 system roles
		t.Errorf("GetAllRoles failed: Should have fetched all items but got: %v", len(response))
	}
}
