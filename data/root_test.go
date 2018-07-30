package data_test

import (
	"fmt"
	"os"
)

//	Gets the database path for this environment:
func getTestFile() string {
	return fmt.Sprintf("%s/testdatabase.db", os.TempDir())
}
