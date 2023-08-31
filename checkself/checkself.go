package checkself

import (
	"context"
	"fmt"
	"os"
	"strings"

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

func printTitle(title string) {
	fmt.Printf("\n-- %s %s\n\n", title, strings.Repeat("-", 76-len(title)))
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
	reportInfo := []string{}
	reportErrors := []string{}

	// List tests.  Order is important.

	testFunctions := []func(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error){
		checkself.Prolog,
		checkself.ListEnvironmentVariables,
		checkself.ListStructVariables,
		checkself.CheckConfigPath,
		checkself.CheckResourcePath,
		checkself.CheckSupportPath,
		checkself.CheckDatabaseUrl,
		checkself.CheckEngineConfigurationJson,
		checkself.Break,
	}

	// Perform checks.

	for _, testFunction := range testFunctions {
		reportChecks, reportInfo, reportErrors, err = testFunction(ctx, reportChecks, reportInfo, reportErrors)
		if err != nil {
			if len(err.Error()) > 0 {
				reportErrors = append(reportErrors, err.Error())
			}
			break
		}
	}

	// Output reports.

	if len(reportInfo) > 0 {
		printTitle("Information")
		for _, message := range reportInfo {
			fmt.Println(message)
		}
	}

	if len(reportChecks) > 0 {
		printTitle("Checks performed")
		for index, message := range reportChecks {
			fmt.Printf("%6d. %s\n", index+1, message)
		}
	}

	if len(reportErrors) > 0 {
		printTitle("Errors")
		for index, message := range reportErrors {
			fmt.Printf("%6d. %s\n\n", index+1, message)
		}
		err = fmt.Errorf("%d errors detected", len(reportErrors))
		fmt.Printf("Result: %s\n", err.Error())
	} else {
		printTitle("Result")
		fmt.Printf("No errors detected.\n")
	}
	fmt.Printf("%s\n\n", strings.Repeat("-", 80))

	return err
}
