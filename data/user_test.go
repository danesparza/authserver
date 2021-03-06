package data_test

import (
	"os"
	"testing"

	"github.com/danesparza/authserver/data"
)

func TestUser_Database_ShouldNotExistYet(t *testing.T) {
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

func TestUser_AddUser_Successful(t *testing.T) {
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
	response, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := response

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
		t.Errorf("AddUser failed: Should have added an item without error: %s", errAdd)
	}

	if newUser.Created.IsZero() || newUser.Updated.IsZero() {
		t.Errorf("AddUser failed: Should have set an item with the correct datetime: %+v", newUser)
	}

	if newUser.CreatedBy != uctx.Name {
		t.Errorf("AddUser failed: Should have set an item with the correct 'created by' user: %+v", newUser)
	}

	if newUser.UpdatedBy != uctx.Name {
		t.Errorf("AddUser failed: Should have set an item with the correct 'updated by' user: %+v", newUser)
	}
}

func TestUser_AddDuplicateUser_ReturnsError(t *testing.T) {
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
	response, _, err := db.AuthSystemBootstrap()
	if err != nil {
		t.Errorf("AuthSystemBootstrap failed: Should have bootstrapped without error: %s", err)
	}

	//	Our 'context' user (the one performing the action)
	uctx := response

	//	Create new user:
	u1 := data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret",
		Description: "Unit test user",
	}

	//	Create duplicate user:
	u2 := data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret",
		Description: "Unit test duplicate user",
	}

	//	Our new password:
	userPassword := "newpassword"

	//	Act
	_, errAdd := db.AddUser(uctx, u1, userPassword)
	_, errAdd2 := db.AddUser(uctx, u2, userPassword)

	//	Assert
	if errAdd != nil {
		t.Errorf("AddUser failed: Should have added an item without error: %s", errAdd)
	}

	if errAdd2 == nil {
		t.Errorf("AddUser failed: Should have returned error (duplicate user name), but didn't")
	}

}

func TestUser_GetAllUsers_NoItems_NoErrors(t *testing.T) {
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

	//	Try storing some users:
	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser1",
		SecretHash:  "SomeRandomSecret1",
		Description: "Unit test user 1",
	}, userPassword)
	if err != nil {
		t.Errorf("AddUser 1 failed: Should have created users without error: %s", err)
	}

	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	}, userPassword)
	if err != nil {
		t.Errorf("AddUser 2 failed: Should have created users without error: %s", err)
	}

	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser3",
		SecretHash:  "SomeRandomSecret3",
		Description: "Unit test user 3",
	}, userPassword)
	if err != nil {
		t.Errorf("AddUser 3 failed: Should have created users without error: %s", err)
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

func TestUser_AddUser_NoCredentials_ReturnsError(t *testing.T) {
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

	uctx = newUser1

	//	Act
	//	Make the context user that user and try to add a user:
	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	}, userPassword)

	//	Assert
	if err == nil {
		t.Errorf("AddUser failed: Should not have added user2 because context user didn't have permission")
	}
}

func TestUser_AddUser_ValidCredentials_Successful(t *testing.T) {
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

	//	Act
	//	Add user to the system / admin resource & role
	_, err = db.AddUserToResourceWithRole(uctx, newUser1, data.Resource{ID: "bdldpjad2pm0cd64ra81"}, data.Role{ID: "bdldpjad2pm0cd64ra82"})
	if err != nil {
		t.Errorf("AddUserToResourceWithRole failed: Should have added newUser1 to the system admin resource and role without error, but got: %s", err)
	}

	//	Make the new user the context user
	uctx = newUser1

	//	Try to add a user with the newly created system admin:
	_, err = db.AddUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	}, userPassword)

	//	Assert
	if err != nil {
		t.Errorf("AddUser failed: Should have created a new user without an error, but got: %s", err)
	}
}

func TestUser_AddUser_NewResourceAndRole_Successful(t *testing.T) {
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

	//	Create a new resource
	newResource1, err := db.AddResource(uctx, data.Resource{Name: "Application 1", Description: "Unit test application"})
	if err != nil {
		t.Errorf("AddResource failed: Should have created newResource1 without issue, but got error: %s", err)
	}

	//	Add the user to the resource as a resource delegate
	_, err = db.AddUserToResourceWithRole(uctx, newUser1, data.Resource{ID: newResource1.ID}, data.Role{ID: data.BuiltIn.ResourceDelegateRole})
	if err != nil {
		t.Errorf("AddUserToResourceWithRole failed: Should have added newUser1 to the system admin resource and role without error, but got: %s", err)
	}

	//	Act / Assert

	//	Make the new user the context user
	uctx = newUser1

	//	Try to add a user with the newly created resource delgate:
	newAppUser, err := db.AddUser(uctx, data.User{
		Name:        "TestUser2",
		SecretHash:  "SomeRandomSecret2",
		Description: "Unit test user 2",
	}, userPassword)

	if err != nil {
		t.Errorf("AddUser failed: Should have created a new user (as a resource delegate) without an error, but got: %s", err)
	}

	//	Add a new role
	newAppRole, err := db.AddRole(uctx, data.Role{
		Name: "Special_access",
	})

	if err != nil {
		t.Errorf("AddRole failed: Should have created a new role (as a resource delegate) without an error, but got: %s", err)
	}

	//	Add the new user to the new resource and role (with the resource delegate as the context user)
	_, err = db.AddUserToResourceWithRole(uctx, newAppUser, newResource1, newAppRole)
	if err != nil {
		t.Errorf("AddUserToResourceWithRole failed: Should have added newAppUser to the newResource1 and newAppRole without error, but got: %s", err)
	}
}
