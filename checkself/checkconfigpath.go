package checkself

import (
	"context"
	"fmt"
	"os"

	"github.com/senzing/go-cmdhelping/option"
)

func (checkself *CheckSelfImpl) CheckConfigPath(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.ConfigPath) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Config path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s=%s", option.ConfigPath.Envar, checkself.ConfigPath))
	requiredFiles := []string{
		"cfgVariant.json",
		"defaultGNRCP.config",
	}
	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", checkself.ConfigPath, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", option.ConfigPath.Envar, checkself.ConfigPath, targetFile))
		}
	}

	return reportChecks, reportErrors, err
}
