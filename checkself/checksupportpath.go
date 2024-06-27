package checkself

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-cmdhelping/option"
)

var RequiredSupportFiles = []string{
	"anyTransRule.ibm",
	"g2SifterRules.ibm",
}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckSupportPath(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	_ = ctx

	// Short-circuit exit.

	if len(checkself.SupportPath) == 0 {
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Prolog.

	reportChecks = append(reportChecks, fmt.Sprintf("Check support path: %s = %s", option.SupportPath.Envar, checkself.SupportPath))

	// Check Resource path.

	errorList := statFiles(option.SupportPath.Envar, checkself.SupportPath, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
