//go:build linux

package checkself

import (
	"context"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestCheckSelfImpl_CheckSelf_Null(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

func TestCheckSelfImpl_CheckSelf_EngineConfigurationJson(test *testing.T) {
	ctx := context.TODO()
	engineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	assert.Nil(test, err)
	testObject := &CheckSelfImpl{
		EngineConfigurationJson: engineConfigurationJson,
	}
	err = testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleCheckSelfImpl_CheckSelf() {
	// For more information, visit https://github.com/Senzing/check-self/blob/main/examplepackage/examplepackage_test.go
	ctx := context.TODO()
	examplePackage := &CheckSelfImpl{}
	examplePackage.CheckSelf(ctx)
	//Output:
	// Checks performed:

	// Done. No errors detected.
}
