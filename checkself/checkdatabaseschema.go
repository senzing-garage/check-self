package checkself

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-databasing/checker"
	"github.com/senzing-garage/go-databasing/connector"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckDatabaseSchema(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.DatabaseURL) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check database schema for %s", checkself.DatabaseURL))

	// Connect to the database.

	databaseConnector, err := connector.NewConnector(ctx, checkself.DatabaseURL)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not create a database connector. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseURL.Envar, checkself.DatabaseURL, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Check for Senzing database schema.

	checker := &checker.BasicChecker{
		DatabaseConnector: databaseConnector,
	}
	isSchemaInstalled, err := checker.IsSchemaInstalled(ctx)
	if !isSchemaInstalled {
		reportErrors = append(reportErrors, fmt.Sprintf("Senzing database schema has not been installed in %s. For more information, visit https://hub.senzing.com/...  Error: %s", checkself.DatabaseURL, err.Error()))
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
