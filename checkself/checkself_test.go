package checkself

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestCheckSelfImpl_CheckSelf_Null(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{}
	err := testObject.CheckSelf(ctx)
	assert.Nil(test, err)
}

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
