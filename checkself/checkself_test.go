package checkself

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-common/g2engineconfigurationjson"
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
	engineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	assert.Nil(test, err)
	testObject := &CheckSelfImpl{
		EngineConfigurationJson: engineConfigurationJson,
	}
	err = testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
