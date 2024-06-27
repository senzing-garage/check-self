package checkself

import (
	"context"
	"testing"

	"github.com/senzing-garage/go-helpers/settings"
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

func TestCheckSelfImpl_CheckSelf_settings(test *testing.T) {
	ctx := context.TODO()
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	assert.Nil(test, err)
	testObject := &CheckSelfImpl{
		Settings: settings,
	}
	err = testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}
