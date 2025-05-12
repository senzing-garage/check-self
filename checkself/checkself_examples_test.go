//go:build linux

package checkself

import (
	"context"
	"fmt"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleBasicCheckSelf_CheckSelf() {
	// For more information, visit https://github.com/senzing-garage/check-self/blob/main/checkself/checkself_examples_test.go
	ctx := context.TODO()
	examplePackage := &BasicCheckSelf{}

	err := examplePackage.CheckSelf(ctx)
	if err != nil {
		fmt.Print(err)
	}
}
