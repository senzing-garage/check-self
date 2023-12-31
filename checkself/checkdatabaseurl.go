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

func checkDatabaseUrl(ctx context.Context, variableName string, databaseUrl string) []string {
	reportErrors := []string{}

	// Parse the database URL.

	parsedUrl, err := url.Parse(databaseUrl)
	if err != nil {
		if strings.HasPrefix(databaseUrl, "postgresql") {
			index := strings.LastIndex(databaseUrl, ":")
			newDatabaseUrl := databaseUrl[:index] + "/" + databaseUrl[index+1:]
			parsedUrl, err = url.Parse(newDatabaseUrl)
		}
		if err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not parse database URL. For more information, visit https://hub.senzing.com/...  Error: %s", variableName, databaseUrl, err.Error()))
			return reportErrors
		}
	}

	// Check database URL scheme.

	if len(parsedUrl.Scheme) == 0 {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. A database scheme is needed (e.g. postgresql://...). For more information, visit https://hub.senzing.com/...", variableName, databaseUrl))
		return reportErrors
	}

	databaseSchemes := []string{
		"sqlite3",
		"postgresql",
		"mysql",
		"mssql",
	}

	if !slices.Contains(databaseSchemes, parsedUrl.Scheme) {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Scheme '%s://' is not recognized. For more information, visit https://hub.senzing.com/...", variableName, databaseUrl, parsedUrl.Scheme))
		return reportErrors
	}

	// Specific database URL scheme checks.

	if parsedUrl.Scheme == "sqlite3" {
		sqliteFilename, err := dbhelper.ExtractSqliteDatabaseFilename(databaseUrl)
		if err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Error: %s. For more information, visit https://hub.senzing.com/...", variableName, databaseUrl, err.Error()))
			return reportErrors
		}
		if _, err := os.Stat(sqliteFilename); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", variableName, databaseUrl, sqliteFilename))
			return reportErrors
		}
	}

	// Check database connector creation.

	databaseConnector, err := connector.NewConnector(ctx, databaseUrl)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not make a new connector. For more information, visit https://hub.senzing.com/...  Error: %s", variableName, databaseUrl, err.Error()))
		return reportErrors
	}

	// Check database connection.

	database := sql.OpenDB(databaseConnector)
	defer database.Close()

	err = database.PingContext(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not connect. For more information, visit https://hub.senzing.com/...  Error: %s", variableName, databaseUrl, err.Error()))
	}

	return reportErrors

}

func (checkself *CheckSelfImpl) CheckDatabaseUrl(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.DatabaseUrl) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check database URL: %s = %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl))

	// Check database URL.

	errorList := checkDatabaseUrl(ctx, option.DatabaseUrl.Envar, checkself.DatabaseUrl)
	reportErrors = append(reportErrors, errorList...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
