package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
)

func (checkself *CheckSelfImpl) CheckDatabaseUrl(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.DatabaseUrl) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Config path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s=%s", option.DatabaseUrl.Envar, checkself.DatabaseUrl))

	return reportChecks, reportErrors, err
}
