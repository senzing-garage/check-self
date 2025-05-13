//go:build darwin

package checkself_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/senzing-garage/check-self/checkself"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicCheckSelf_CheckSelf_Paths(test *testing.T) {
	ctx := context.TODO()
	senzingPath := settings.GetSenzingPath()
	testObject := &checkself.BasicCheckSelf{
		ConfigPath:   fmt.Sprintf("%s/er/etc", senzingPath),
		DatabaseURL:  "sqlite3://na:na@/tmp/sqlite/G2C.db",
		ResourcePath: fmt.Sprintf("%s/er/resources", senzingPath),
		SupportPath:  fmt.Sprintf("%s/data", senzingPath),
	}
	err := testObject.CheckSelf(ctx)
	require.NoError(test, err)
}
