//go:build linux

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
	testObject := &CheckSelfImpl{
		ConfigPath:   "/etc/opt/senzing",
		DatabaseUrl:  "sqlite3://na:na@/tmp/sqlite/G2C.db",
		ResourcePath: "/opt/senzing/g2/resources",
		SupportPath:  "/opt/senzing/data",
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
