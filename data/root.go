package data

// SystemDB is the BoltDB database for
// user/application/role storage
type SystemDB struct {
	Database string
}

// TokenDB is the BoltDB database for
// token storage
type TokenDB struct {
	Database string
}
