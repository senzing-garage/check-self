package checkself

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-databasing/connector"
	"github.com/senzing-garage/go-databasing/dbhelper"
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
	var err error

	// Short-circuit exit.

	if len(checkself.DatabaseURL) == 0 {
		return reportChecks, reportInfo, reportErrors, err
	}

	// Prolog.

	reportChecks = append(
		reportChecks,
		fmt.Sprintf("Check database URL: %s = %s", option.DatabaseURL.Envar, checkself.DatabaseURL),
	)

	// Check database URL.

	reportErrors = append(reportErrors, checkDatabaseURL(ctx, option.DatabaseURL.Envar, checkself.DatabaseURL)...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, err
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func checkDatabaseURL(ctx context.Context, variableName string, databaseURL string) []string {
	result := []string{}

	// Parse the database URL.

	normalizedDatabaseURL := databaseURL
	if strings.HasPrefix(databaseURL, "postgresql") {
		index := strings.LastIndex(databaseURL, ":")
		normalizedDatabaseURL = databaseURL[:index] + "/" + databaseURL[index+1:]
	}

	parsedURL, err := url.Parse(normalizedDatabaseURL)
	if err != nil {
		return append(result, fmt.Sprintf(
			"%s = %s is misconfigured. Could not parse database URL. For more information, visit https://hub.senzing.com/...  Error: %s",
			variableName,
			databaseURL,
			err.Error(),
		))
	}

	// Check database URL scheme.

	if len(parsedURL.Scheme) == 0 {
		return append(result, fmt.Sprintf(
			"%s = %s is misconfigured. A database scheme is needed (e.g. postgresql://...). For more information, visit https://hub.senzing.com/...",
			variableName,
			databaseURL,
		))
	}

	databaseSchemes := []string{
		"sqlite3",
		"postgresql",
		"mysql",
		"mssql",
	}

	if !slices.Contains(databaseSchemes, parsedURL.Scheme) {
		return append(result, fmt.Sprintf(
			"%s = %s is misconfigured. Scheme '%s://' is not recognized. For more information, visit https://hub.senzing.com/...",
			variableName,
			databaseURL,
			parsedURL.Scheme,
		))
	}

	// Specific database URL scheme checks.

	if parsedURL.Scheme == "sqlite3" {
		result = append(result, checkSqlite(variableName, databaseURL)...)
	}

	result = append(result, checkDatabaseConnection(ctx, variableName, databaseURL)...)

	// Check database connector creation.

	return result

}

func checkDatabaseConnection(ctx context.Context, variableName string, databaseURL string) []string {
	var result []string

	databaseConnector, err := connector.NewConnector(ctx, databaseURL)
	if err != nil {
		return append(
			result,
			fmt.Sprintf(
				"%s = %s is misconfigured. Could not make a new connector. For more information, visit https://hub.senzing.com/...  Error: %s",
				variableName,
				databaseURL,
				err.Error(),
			),
		)
	}

	// Check database connection.

	database := sql.OpenDB(databaseConnector)
	defer database.Close()

	err = database.PingContext(ctx)
	if err != nil {
		return append(
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

func checkSqlite(variableName string, databaseURL string) []string {
	var result []string

	sqliteFilename, err := dbhelper.ExtractSqliteDatabaseFilename(databaseURL)
	if err != nil {
		return append(result, fmt.Sprintf(
			"%s = %s is misconfigured. Error: %s. For more information, visit https://hub.senzing.com/...",
			variableName,
			databaseURL,
			err.Error()))
	}
	if _, err := os.Stat(sqliteFilename); err != nil {
		return append(result, fmt.Sprintf(
			"%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...",
			variableName,
			databaseURL,
			sqliteFilename))
	}
	return result
}
