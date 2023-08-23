package checkself

import (
	"context"
	"fmt"
	"os"

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
// Internal methods
// ----------------------------------------------------------------------------

func statFiles(variableName string, path string, requiredFiles []string) []string {
	reportErrors := []string{}
	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", path, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", variableName, path, targetFile))
		}
	}
	return reportErrors
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

	// List tests.  Order is important.

	testFunctions := []func(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error){
		checkself.CheckConfigPath,
		checkself.CheckResourcePath,
		checkself.CheckSupportPath,
		checkself.CheckDatabaseUrl,
		checkself.CheckEngineConfigurationJson,
	}

	// Perform tests.

	for _, testFunction := range testFunctions {
		reportChecks, reportErrors, err = testFunction(ctx, reportChecks, reportErrors)
		if err != nil {
			return err
		}
		if len(reportErrors) > 0 {
			break
		}
	}

	// Output reports.

	fmt.Printf("\nChecks performed:\n\n")
	for index, message := range reportChecks {
		fmt.Printf("  %4d - %s\n", index+1, message)
	}

	if len(reportErrors) > 0 {
		err = fmt.Errorf("%d errors detected", len(reportErrors))
		fmt.Printf("\nErrors: %s:\n\n", err.Error())
		for index, message := range reportErrors {
			fmt.Printf("  %4d - %s\n\n", index+1, message)
		}
	} else {
		fmt.Printf("\n\nDone. No errors detected.\n")
	}

	return err
}
