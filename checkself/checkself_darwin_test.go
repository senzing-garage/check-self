//go:build darwin

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
	testObject := &BasicCheckSelf{
		ConfigPath:   "/opt/senzing/g2/etc",
		DatabaseURL:  "sqlite3://na:na@/tmp/sqlite/G2C.db",
		ResourcePath: "/opt/senzing/g2/resources",
		SupportPath:  "/opt/senzing/g2/data",
	}
	err := testObject.CheckSelf(ctx)
	require.NoError(test, err)
}
