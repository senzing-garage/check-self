package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-databasing/checker"
	"github.com/senzing/go-databasing/connector"
)

func (checkself *CheckSelfImpl) CheckDatabaseSchema(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.DatabaseUrl) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check database schema for %s", checkself.DatabaseUrl))

	// Connect to the database.

	databaseConnector, err := connector.NewConnector(ctx, checkself.DatabaseUrl)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not create a database connector. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Check for Senzing database schema.

	checker := &checker.CheckerImpl{
		DatabaseConnector: databaseConnector,
	}
	isSchemaInstalled, err := checker.IsSchemaInstalled(ctx)
	if !isSchemaInstalled {
		reportErrors = append(reportErrors, fmt.Sprintf("Senzing database schema has not been installed in %s. For more information, visit https://hub.senzing.com/...  Error: %s", checkself.DatabaseUrl, err.Error()))
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
