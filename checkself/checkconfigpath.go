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

	// Check Config path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.ConfigPath.Envar, checkself.ConfigPath))
	errorList := statFiles(option.ConfigPath.Envar, checkself.ConfigPath, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)
	return reportChecks, reportInfo, reportErrors, nil
}
