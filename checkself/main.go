package checkself

import (
	"context"
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
const ExampleConstant = 1
const DefaultSenzingToolsLicenseDaysLeft = "30"
const DefaultSenzingToolsLicenseRecordsPercent = "90"

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// An example variable.
var ExampleVariable = map[int]string{
	1: "Just a string",
}
