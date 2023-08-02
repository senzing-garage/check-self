package checkself

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// CheckSelfImpl is the basic checker.
type CheckSelfImpl struct {
	ConfigPath              string
	DatabaseUrl             string            // TODO:
	EngineConfigurationJson string            // TODO:
	EngineLogLevel          string            // TODO:
	GrpcDialOptions         []grpc.DialOption // TODO:
	GrpcUrl                 string            // TODO:
	InputUrl                string            // TODO:
	LicenseStringBase64     string            // TODO:
	LogLevel                string            // TODO:
	ObserverUrl             string            // TODO:
	ResourcePath            string
	SenzingDirectory        string // TODO:
	SupportPath             string
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The CheckSelf method prints the output of numerous checks.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/
func (checkself *CheckSelfImpl) CheckSelf(ctx context.Context) error {
	var err error = nil

	reportChecks := []string{}
	reportErrors := []string{}

	// Perform tests.

	reportChecks, reportErrors, err = checkself.CheckConfigPath(ctx, reportChecks, reportErrors)
	if err != nil {
		return err
	}

	reportChecks, reportErrors, err = checkself.CheckResourcePath(ctx, reportChecks, reportErrors)
	if err != nil {
		return err
	}

	reportChecks, reportErrors, err = checkself.CheckSupportPath(ctx, reportChecks, reportErrors)
	if err != nil {
		return err
	}

	reportChecks, reportErrors, err = checkself.CheckDatabaseUrl(ctx, reportChecks, reportErrors)
	if err != nil {
		return err
	}

	// Output reports.

	fmt.Printf("\nChecks performed:\n\n")
	for index, message := range reportChecks {
		fmt.Printf("  %4d - %s\n", index+1, message)
	}

	if len(reportErrors) > 0 {
		fmt.Printf("\nErrors: %d errors detected:\n\n", len(reportErrors))
		for index, message := range reportErrors {
			fmt.Printf("  %4d - %s\n\n", index+1, message)
		}
	} else {
		fmt.Printf("\n\nDone. No errors detected.\n")
	}

	return nil
}
