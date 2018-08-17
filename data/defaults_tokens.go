package data

/* Tables */
// tokenSchema defines the schema for the token table
var tokenSchema = `
CREATE TABLE IF NOT EXISTS tokens (
	token string NOT NULL,
	userid string NOT NULL,
	created time NOT NULL,
	createdby string NOT NULL,
	expires time NOT NULL,
	deleted time,
	deletedby string
);`

/* Indices */
var tokenIXToken = `
CREATE UNIQUE INDEX IF NOT EXISTS Token ON tokens (token)`
