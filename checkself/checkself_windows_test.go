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
		DatabaseUrl:  `sqlite3://C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

func TestCheckSelfImpl_CheckSelf_Paths4(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  `sqlite3://nowhere.com:0/C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

func TestCheckSelfImpl_CheckSelf_Paths1(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  `sqlite3://na:na@nowhere.com:0/C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

func TestCheckSelfImpl_CheckSelf_Paths2(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  `sqlite3://na:na@nowhere:0/C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

func TestCheckSelfImpl_CheckSelf_Paths3(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{
		ConfigPath:   `C:\Program Files\Senzing\g2\etc`,
		DatabaseUrl:  `sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db`,
		ResourcePath: `C:\Program Files\Senzing\g2\resources`,
		SupportPath:  `C:\Program Files\Senzing\g2\data`,
	}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
