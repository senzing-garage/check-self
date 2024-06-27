//go:build windows

package checkself

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicCheckSelf_CheckSelf_Paths(test *testing.T) {
	ctx := context.TODO()
	databaseUrl := "sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db"
	test.Logf(">>>>> test: %s\n", databaseUrl)

	test.Log("sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db")

	test.Log(`sqlite3://na:na@nowhere/C:\\Temp\\sqlite\\G2C.db`)

	test.Log(`sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db`)

	testObject := &BasicCheckSelf{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseURL:  databaseUrl,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	require.NoError(test, err)
}
