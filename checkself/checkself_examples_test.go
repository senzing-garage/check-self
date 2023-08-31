//go:build linux

package checkself

import (
	"context"
)

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
