package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
)

var RequiredResourceFiles = []string{
	"templates/g2config.json",
}

func (checkself *CheckSelfImpl) CheckResourcePath(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.ResourcePath) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Check Resource path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.ResourcePath.Envar, checkself.ResourcePath))
	errorList := statFiles(option.ResourcePath.Envar, checkself.ResourcePath, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)
	return reportChecks, reportInfo, reportErrors, nil
}
