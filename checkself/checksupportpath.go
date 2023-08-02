package checkself

import (
	"context"
	"fmt"
	"os"

	"github.com/senzing/go-cmdhelping/option"
)

func (checkself *CheckSelfImpl) CheckSupportPath(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.SupportPath) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Resource path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s=%s", option.SupportPath.Envar, checkself.SupportPath))
	requiredFiles := []string{
		"anyTransRule.ibm",
		"g2SifterRules.ibm",
	}
	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", checkself.SupportPath, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", option.SupportPath.Envar, checkself.SupportPath, targetFile))
		}
	}

	return reportChecks, reportErrors, err
}
