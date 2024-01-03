//go:build linux

package checkself

import (
	"context"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleCheckSelfImpl_CheckSelf() {
	// For more information, visit https://github.com/senzing-garage/check-self/blob/main/checkself/checkself_examples_test.go
	ctx := context.TODO()
	examplePackage := &CheckSelfImpl{}
	examplePackage.CheckSelf(ctx)
	//Output:
	// Checks performed:

	// Done. No errors detected.
}
