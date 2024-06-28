package checkself

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-cmdhelping/option"
)

var RequiredResourceFiles = []string{
	"templates/g2config.json",
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckResourcePath(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	_ = ctx

	// Short-circuit exit.

	if len(checkself.ResourcePath) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check resource path: %s = %s", option.ResourcePath.Envar, checkself.ResourcePath))

	// Check Resource path.

	errorList := statFiles(option.ResourcePath.Envar, checkself.ResourcePath, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
