package checkself

import (
	"context"
	"errors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The ExamplePackage interface is an example interface.
type CheckSelf interface {
	CheckSelf(ctx context.Context) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// An example constant.
const (
	ExampleConstant                          = 1
	DefaultSenzingToolsLicenseDaysLeft       = "30"
	DefaultSenzingToolsLicenseRecordsPercent = "90"
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// An example variable.
var ExampleVariable = map[int]string{
	1: "Just a string",
}

var errForPackage = errors.New("checkself")
