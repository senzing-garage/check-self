package checkself

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-databasing/connector"
	"golang.org/x/exp/slices"
)

func (checkself *CheckSelfImpl) CheckDatabaseUrl(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.DatabaseUrl) == 0 {
		return reportChecks, reportErrors, err
	}

	reportChecks = append(reportChecks, fmt.Sprintf("%s=%s", option.DatabaseUrl.Envar, checkself.DatabaseUrl))

	// Parse the database URL.

	parsedUrl, err := url.Parse(checkself.DatabaseUrl)
	if err != nil {
		if strings.HasPrefix(checkself.DatabaseUrl, "postgresql") {
			index := strings.LastIndex(checkself.DatabaseUrl, ":")
			newDatabaseUrl := checkself.DatabaseUrl[:index] + "/" + checkself.DatabaseUrl[index+1:]
			parsedUrl, err = url.Parse(newDatabaseUrl)
		}
		if err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not parse database URL. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl, err.Error()))
			return reportChecks, reportErrors, err
		}
	}

	// Check database URL scheme.

	if len(parsedUrl.Scheme) == 0 {
		reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. A database scheme is needed (e.g. postgresql://...). For more information, visit https://hub.senzing.com/...", option.DatabaseUrl.Envar, checkself.DatabaseUrl))
		return reportChecks, reportErrors, err
	}

	databaseSchemes := []string{
		"sqlite3",
		"postgresql",
		"mysql",
		"mssql",
	}

	if !slices.Contains(databaseSchemes, parsedUrl.Scheme) {
		reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Scheme '%s://' is not recognized. For more information, visit https://hub.senzing.com/...", option.DatabaseUrl.Envar, checkself.DatabaseUrl, parsedUrl.Scheme))
		return reportChecks, reportErrors, err
	}

	// Check database connector creation.

	databaseConnector, err := connector.NewConnector(ctx, checkself.DatabaseUrl)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not make a new connector. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl, err.Error()))
		return reportChecks, reportErrors, err
	}

	// Check database connection.

	database := sql.OpenDB(databaseConnector)
	defer database.Close()

	err = database.PingContext(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not connect. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl, err.Error()))
	}

	return reportChecks, reportErrors, nil
}
