package checkself

import (
	"context"
	"fmt"
	"os"
	"testing"

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

func TestCheckSelfImpl_CheckSelf(test *testing.T) {
	ctx := context.TODO()
	testObject := &CheckSelfImpl{}
	err := testObject.CheckSelf(ctx)
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
