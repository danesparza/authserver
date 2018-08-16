package data

import (
	"database/sql"
	"fmt"

	// QL sql driver
	_ "github.com/cznic/ql/driver"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// SystemDB is the BoltDB database for
// user/application/role storage
type SystemDB struct {
	db *sql.DB
}

// TokenDB is the BoltDB database for
// token storage
type TokenDB struct {
	db *sql.DB
}

// NewSystemDB creates a new instance of a SystemDB
func NewSystemDB(dbpath string) (*SystemDB, error) {
	retval := new(SystemDB)

	//	Create a reference to our bolt db
	db, err := sql.Open("ql", dbpath)
	if err != nil {
		return nil, fmt.Errorf("An error occurred opening the SystemDB: %s", err)
	}
	retval.db = db

	//	Return our systemdb reference
	return retval, nil
}

// Close closes the SystemDB database
func (store SystemDB) Close() error {
	return store.db.Close()
}

// AuthSystemBootstrap initializes the SystemDB and creates any default admin users / roles / resources
func (store SystemDB) AuthSystemBootstrap() (User, string, error) {
	adminUser := User{}
	adminPassword := ""

	//	Start our database transaction
	tx, err := store.db.Begin()
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem starting a transaction to bootstrap auth system")
	}

	//	Create our database schema and indices
	//	Resource schema / indices
	_, err = tx.Exec(resourceSchema)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding resource schema: %s", err)
	}
	_, err = tx.Exec(resourceIXSysID)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding resource id index: %s", err)
	}
	_, err = tx.Exec(resourceIXName)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding resource name index: %s", err)
	}

	//	Role schema / indices
	_, err = tx.Exec(roleSchema)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding role schema: %s", err)
	}
	_, err = tx.Exec(roleIXSysID)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding role id index: %s", err)
	}

	//	User schema / indices
	_, err = tx.Exec(userSchema)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding user schema: %s", err)
	}
	_, err = tx.Exec(userIXSysID)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding user id index: %s", err)
	}
	_, err = tx.Exec(userIXName)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding user name index: %s", err)
	}

	//	UserResourceRole schema / indices
	_, err = tx.Exec(userResourceRoleSchema)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding user_resource_role schema: %s", err)
	}
	_, err = tx.Exec(userResourceRoleIXID)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding user_resource_role id index: %s", err)
	}

	//	Generate a password for the admin user
	adminPassword = xid.New().String()

	//	Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem hashing admin password: %s", err)
	}

	//	Add our default admin user - the insert statement requires some parameters be passed:
	_, err = tx.Exec(defaultAdminUser, BuiltIn.AdminUser, string(hashedPassword))
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding admin user: %s", err)
	}

	//	Create the default system resources:
	_, err = tx.Exec(defaultSystemResource, BuiltIn.SystemResource)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding system resource: %s", err)
	}

	//	Create the default system roles (id/name/description):
	_, err = tx.Exec(defaultSystemRole, BuiltIn.AdminRole, "sys_admin", "System admin role")
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding system role: %s", err)
	}

	_, err = tx.Exec(defaultSystemRole, BuiltIn.ResourceDelegateRole, "sys_delegate", "Resource delegate role")
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding system role: %s", err)
	}

	//	Create the default system credentials (user/resource/role):
	//	-- Admins get the sysadmin role
	_, err = tx.Exec(defaultSystemCredentials, BuiltIn.AdminUser, BuiltIn.SystemResource, BuiltIn.AdminRole)
	if err != nil {
		tx.Rollback()
		return adminUser, adminPassword, fmt.Errorf("Problem adding system credential: %s", err)
	}

	//	Commit our transaction
	err = tx.Commit()
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem committing a transaction to bootstrap auth system")
	}

	//	Get our admin user from the database and create our return object:
	adminUser = User{}
	err = store.db.QueryRow("SELECT id, enabled, name, description, secrethash, created, createdby, updated, updatedby, deleted, deletedby FROM user WHERE id=$1;", BuiltIn.AdminUser).Scan(
		&adminUser.ID,
		&adminUser.Enabled,
		&adminUser.Name,
		&adminUser.Description,
		&adminUser.SecretHash,
		&adminUser.Created,
		&adminUser.CreatedBy,
		&adminUser.Updated,
		&adminUser.UpdatedBy,
		&adminUser.Deleted,
		&adminUser.DeletedBy)
	if err != nil {
		return adminUser, adminPassword, fmt.Errorf("Problem selecting admin user: %s", err)
	}

	return adminUser, adminPassword, nil
}
