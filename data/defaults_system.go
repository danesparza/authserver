package data

// Defaults encapsulates system defaults
type Defaults struct {
	// AdminUser is the admin user id
	AdminUser string

	// SystemResource is the system resource id
	SystemResource string

	// AdminRole is the system admin role id
	AdminRole string

	// ResourceDelegateRole is the resource delegate role id
	ResourceDelegateRole string
}

// BuiltIn is a catalog of system default values
var BuiltIn Defaults

/* Tables */
// resourceSchema defines the schema for the resource table
var resourceSchema = `
CREATE TABLE IF NOT EXISTS resource (
	id string NOT NULL,
	name string NOT NULL,
    description string,
	created time NOT NULL,
	createdby string NOT NULL,
	updated time NOT NULL,
	updatedby string NOT NULL,
	deleted time,
	deletedby string
);`

// roleSchema defines the schema for the role table
var roleSchema = `
CREATE TABLE IF NOT EXISTS role (
	id string NOT NULL,
	name string NOT NULL,
    description string,
	created time NOT NULL,
	createdby string NOT NULL,
	updated time NOT NULL,
	updatedby string NOT NULL,
	deleted time,
	deletedby string
);`

// userSchema defines the schema for the user table
var userSchema = `
CREATE TABLE IF NOT EXISTS user (
	id string NOT NULL,
	enabled bool NOT NULL,
    name string NOT NULL,
	description string,
	secrethash string,
	created time NOT NULL,
	createdby string NOT NULL,
	updated time NOT NULL,
	updatedby string NOT NULL,
	deleted time,
	deletedby string
);`

// userResourceRoleSchema defines the schema for the user_resource_role table
var userResourceRoleSchema = `
CREATE TABLE IF NOT EXISTS user_resource_role (
	userid string NOT NULL,
	resourceid string NOT NULL,
	roleid string NOT NULL,
	created time NOT NULL,
	createdby string NOT NULL,
	updated time NOT NULL,
	updatedby string NOT NULL,
	deleted time,
	deletedby string
);`

/* Indices */
var userIXSysID = `
CREATE UNIQUE INDEX IF NOT EXISTS UserID ON user (id)`

var userIXName = `
CREATE UNIQUE INDEX IF NOT EXISTS UserName ON user (name)`

var roleIXSysID = `
CREATE UNIQUE INDEX IF NOT EXISTS RoleID ON role (id)`

var resourceIXSysID = `
CREATE UNIQUE INDEX IF NOT EXISTS ResourceID ON resource (id)`

var resourceIXName = `
CREATE UNIQUE INDEX IF NOT EXISTS ResourceName ON resource (name)`

var userResourceRoleIXID = `
CREATE UNIQUE INDEX IF NOT EXISTS UserResourceRoleID ON user_resource_role (userid, resourceid, roleid)`

// defaultAdminUser is the insert statement that creates the default admin user - it requires 2 parameters:
// - the id of the admin user
// - the generated secrethash for the admin user's password
var defaultAdminUser = `
INSERT INTO 
	user(id, enabled, name, description, secrethash, created, createdby, updated, updatedby) 
	values($1, true, "admin", "Default admin user", $2, now(), "system", now(), "system");`

// defaultSystemResource is the insert statement that creates the default system resource - it requires 1 parameter:
// - the id of the system resource
var defaultSystemResource = `
INSERT INTO 
	resource(id, name, description, created, createdby, updated, updatedby) 
	values($1, "system", "Default authsystem resource", now(), "system", now(), "system");`

// defaultAdminUser is the insert statement that creates the default system roles - it requires 3 parameters:
// - the id of the system role
// - the name of the system role
// - the description of the system role
var defaultSystemRole = `
	INSERT INTO 
		role(id, name, description, created, createdby, updated, updatedby) 
		values($1, $2, $3, now(), "system", now(), "system");`

// defaultSystemCredentials is the insert statement that creates the default system credentials - it requires 3 parameters:
// - the id of the admin user
// - the id of the system resource
// - the id of the system role
var defaultSystemCredentials = `
	INSERT INTO
		user_resource_role(userid, resourceid, roleid, created, createdby, updated, updatedby)
		values($1, $2, $3, now(), "system", now(), "system")
`

// getResourcesForUser is the query to get all resources for a given user.  It requires 1 parameter:
// - the id of the user to check
var getResourcesForUser = `
select 
	distinct resource.id, resource.name, resource.description 
from 
	resource, user_resource_role 
where 
	resource.id = user_resource_role.resourceid 
	and user_resource_role.userid = $1
`

// getRolesForUserAndResources is the query to get all roles for a given user and resource.  It requires 2 parameters:
// - the id of the user to check
// - the id of the resource to check
var getRolesForUserAndResources = `
select 
	distinct role.id, role.name, role.description 
from 
	role, user_resource_role 
where 
	role.id = user_resource_role.roleid
	and user_resource_role.userid = $1
	and user_resource_role.resourceid = $2  	
`

func init() {

	//	Setup our defaults
	BuiltIn = Defaults{
		AdminUser:            "bdldpjad2pm0cd64ra80",
		SystemResource:       "bdldpjad2pm0cd64ra81",
		AdminRole:            "bdldpjad2pm0cd64ra82",
		ResourceDelegateRole: "bdldpjad2pm0cd64ra83",
	}
}
