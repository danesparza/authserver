package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestResource_Database_ShouldNotExistYet(t *testing.T) {
	//	Arrange
	filename := getTestFile()

	//	Act

	//	Assert
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("System database file check failed: System file %s already exists, and shouldn't", filename)
	}
}

func TestResource_Add_Successful(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewDBManager(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	uctx, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Create new resource:
	r1 := data.Resource{
		Name:        "TestResource1",
		Description: "Unit test resource",
	}

	//	Act
	response, err := db.AddResource(uctx, r1)

	//	Assert
	if err != nil {
		t.Errorf("SetResource failed: Should have set an item without error: %s", err)
	}

	if response.Created.IsZero() || response.Updated.IsZero() {
		t.Errorf("SetResource failed: Should have set an item with the correct datetime: %+v", response)
	}

	if response.CreatedBy != uctx.Name {
		t.Errorf("SetResource failed: Should have set an item with the correct 'created by' user: %+v", response)
	}

	if response.UpdatedBy != uctx.Name {
		t.Errorf("SetResource failed: Should have set an item with the correct 'updated by' user: %+v", response)
	}
}

func TestResource_GetAllResources_NoItems_NoErrors(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewDBManager(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	uctx, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Act
	response, err := db.GetAllResources(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllResources failed: Should have gotten all items without error: %s", err)
	}

	if len(response) != 1 { // Bootstrap adds a system resource
		t.Errorf("GetAllResources failed: Should have only fetched the default resource")
	}
}

func TestResource_GetAllResources_ItemsInDB_ReturnsItems(t *testing.T) {
	//	Arrange
	filename := getTestFile()
	defer os.Remove(filename)

	db, err := data.NewDBManager(filename)
	if err != nil {
		t.Errorf("NewSystemDB failed: %s", err)
	}
	defer db.Close()

	//	Bootstrap
	uctx, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Try storing some resources:
	_, err = db.AddResource(uctx, data.Resource{
		Name:        "TestResource1",
		Description: "Unit test resource 1",
	})
	if err != nil {
		t.Errorf("SetResource 1 failed: Should have created item without error: %s", err)
	}

	_, err = db.AddResource(uctx, data.Resource{
		Name:        "TestResource2",
		Description: "Unit test resource 2",
	})
	if err != nil {
		t.Errorf("SetResource 2 failed: Should have created item without error: %s", err)
	}

	//	Act
	response, err := db.GetAllResources(uctx)

	//	Assert
	if err != nil {
		t.Errorf("GetAllResources failed: Should have gotten the items without error: %s", err)
	}

	if len(response) != 3 { // Bootstrap adds a system resource
		t.Errorf("GetAllResources failed: Should have gotten all items.  Instead, got: %v", len(response))
	}
}
