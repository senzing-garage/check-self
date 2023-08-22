package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
)

var RequiredResourceFiles = []string{
	"templates/g2config.json",
}

func (checkself *CheckSelfImpl) CheckResourcePath(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.ResourcePath) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Resource path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.ResourcePath.Envar, checkself.ResourcePath))
	errorList := statFiles(option.ResourcePath.Envar, checkself.ResourcePath, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)
	return reportChecks, reportErrors, err
}
