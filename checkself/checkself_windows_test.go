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
	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  `sqlite3://na:na@/C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
