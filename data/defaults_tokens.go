package data

/* Tables */
// tokenSchema defines the schema for the token table
var tokenSchema = `
CREATE TABLE IF NOT EXISTS tokens (
	token string NOT NULL,
	userid string NOT NULL,
	created time NOT NULL,
	expires time NOT NULL,
	deleted time,
	deletedby string
);`

/* Indices */
var tokenIXToken = `
CREATE UNIQUE INDEX IF NOT EXISTS TokenID ON tokens (token)`

var tokenIXUserID = `
CREATE INDEX IF NOT EXISTS TokenUser ON tokens (userid)`
