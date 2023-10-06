//go:build windows

package checkself

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestCheckSelfImpl_CheckSelf_Paths(test *testing.T) {
	ctx := context.TODO()
	databaseUrl := "sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db"
	test.Logf(">>>>> test: %s\n", databaseUrl)

	test.Log("sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db")

	test.Log(`sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db`)

	test.Log(`sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db`)

	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  databaseUrl,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
