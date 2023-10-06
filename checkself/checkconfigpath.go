package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
)

var RequiredConfigFiles = []string{
	"cfgVariant.json",
	"defaultGNRCP.config",
}

func (checkself *CheckSelfImpl) CheckConfigPath(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.ConfigPath) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check configuration path: %s = %s", option.ConfigPath.Envar, checkself.ConfigPath))

	// Check Config path.

	errorList := statFiles(option.ConfigPath.Envar, checkself.ConfigPath, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
