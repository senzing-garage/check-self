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

func (checkself *CheckSelfImpl) CheckSupportPath(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.SupportPath) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Check Resource path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.SupportPath.Envar, checkself.SupportPath))
	errorList := statFiles(option.SupportPath.Envar, checkself.SupportPath, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)
	return reportChecks, reportInfo, reportErrors, nil
}
