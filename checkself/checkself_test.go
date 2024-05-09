package checkself

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestCheckSelfImpl_CheckSelf_Null(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{}
	err := testObject.CheckSelf(ctx)
	assert.NotNil(test, err)
}

func TestCheckSelfImpl_CheckSelf_EngineConfigurationJson(test *testing.T) {
	ctx := context.TODO()
	engineConfigurationJson, err := engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	assert.Nil(test, err)
	testObject := &CheckSelfImpl{
		Settings: engineConfigurationJson,
	}
	err = testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
