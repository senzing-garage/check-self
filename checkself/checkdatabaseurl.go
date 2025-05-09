package checkself

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-databasing/connector"
	"github.com/senzing-garage/go-databasing/dbhelper"
	"golang.org/x/exp/slices"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckDatabaseURL(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.DatabaseURL) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(
		reportChecks,
		fmt.Sprintf("Check database URL: %s = %s", option.DatabaseURL.Envar, checkself.DatabaseURL),
	)

	// Check database URL.

	errorList := checkDatabaseURL(ctx, option.DatabaseURL.Envar, checkself.DatabaseURL)
	reportErrors = append(reportErrors, errorList...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func checkDatabaseURL(ctx context.Context, variableName string, databaseURL string) []string {
	reportErrors := []string{}

	// Parse the database URL.

	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		if strings.HasPrefix(databaseURL, "postgresql") {
			index := strings.LastIndex(databaseURL, ":")
			newDatabaseURL := databaseURL[:index] + "/" + databaseURL[index+1:]
			parsedURL, err = url.Parse(newDatabaseURL)
		}
		if err != nil {
			reportErrors = append(
				reportErrors,
				fmt.Sprintf(
					"%s = %s is misconfigured. Could not parse database URL. For more information, visit https://hub.senzing.com/...  Error: %s",
					variableName,
					databaseURL,
					err.Error(),
				),
			)
			return reportErrors
		}
	}

	// Check database URL scheme.

	if len(parsedURL.Scheme) == 0 {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"%s = %s is misconfigured. A database scheme is needed (e.g. postgresql://...). For more information, visit https://hub.senzing.com/...",
				variableName,
				databaseURL,
			),
		)
		return reportErrors
	}

	databaseSchemes := []string{
		"sqlite3",
		"postgresql",
		"mysql",
		"mssql",
	}

	if !slices.Contains(databaseSchemes, parsedURL.Scheme) {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"%s = %s is misconfigured. Scheme '%s://' is not recognized. For more information, visit https://hub.senzing.com/...",
				variableName,
				databaseURL,
				parsedURL.Scheme,
			),
		)
		return reportErrors
	}

	// Specific database URL scheme checks.

	if parsedURL.Scheme == "sqlite3" {
		sqliteFilename, err := dbhelper.ExtractSqliteDatabaseFilename(databaseURL)
		if err != nil {
			reportErrors = append(
				reportErrors,
				fmt.Sprintf(
					"%s = %s is misconfigured. Error: %s. For more information, visit https://hub.senzing.com/...",
					variableName,
					databaseURL,
					err.Error(),
				),
			)
			return reportErrors
		}
		if _, err := os.Stat(sqliteFilename); err != nil {
			reportErrors = append(
				reportErrors,
				fmt.Sprintf(
					"%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...",
					variableName,
					databaseURL,
					sqliteFilename,
				),
			)
			return reportErrors
		}
	}

	databaseConnectionReport := checkDatabaseConnection(ctx, variableName, databaseURL)
	reportErrors = append(reportErrors, databaseConnectionReport...)

	// Check database connector creation.

	return reportErrors

}

func checkDatabaseConnection(ctx context.Context, variableName string, databaseURL string) []string {
	var result []string

	databaseConnector, err := connector.NewConnector(ctx, databaseURL)
	if err != nil {
		result = append(
			result,
			fmt.Sprintf(
				"%s = %s is misconfigured. Could not make a new connector. For more information, visit https://hub.senzing.com/...  Error: %s",
				variableName,
				databaseURL,
				err.Error(),
			),
		)
		return result
	}

	// Check database connection.

	database := sql.OpenDB(databaseConnector)
	defer database.Close()

	err = database.PingContext(ctx)
	if err != nil {
		result = append(
			result,
			fmt.Sprintf(
				"%s = %s is misconfigured. Could not connect. For more information, visit https://hub.senzing.com/...  Error: %s",
				variableName,
				databaseURL,
				err.Error(),
			),
		)
	}
	return result
}
