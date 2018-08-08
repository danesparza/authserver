package data

// resourceSchema defines the schema for the resource table
var resourceSchema = `
CREATE TABLE IF NOT EXISTS resource (
	sysid string NOT NULL,
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
	sysid string NOT NULL,
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
	sysid string NOT NULL,
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

// adminSysID is the system id of the default admin user
var adminSysID = "bdldpjad2pm0cd64ra80"

// defaultAdmin is the insert statement that creates the default admin user - it requires 2 parameters:
// - the sysid of the admin user
// - the generated secrethash for the admin user's password
var defaultAdmin = `
INSERT INTO 
	user(sysid, enabled, name, description, secrethash, created, createdby, updated, updatedby) 
	values($1, true, "admin", "Default admin user", $2, now(), "system", now(), "system");`
