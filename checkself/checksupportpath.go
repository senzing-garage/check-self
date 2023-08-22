package checkself

import (
	"context"
	"fmt"

	"github.com/senzing/go-cmdhelping/option"
)

var RequiredSupportFiles = []string{
	"anyTransRule.ibm",
	"g2SifterRules.ibm",
}

func (checkself *CheckSelfImpl) CheckSupportPath(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.SupportPath) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Resource path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.SupportPath.Envar, checkself.SupportPath))
	errorList := statFiles(option.SupportPath.Envar, checkself.SupportPath, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)
	return reportChecks, reportErrors, err
}
