package checkself

import (
	"context"
	"fmt"
	"os"

	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-common/engineconfigurationjsonparser"
)

func (checkself *CheckSelfImpl) CheckEngineConfigurationJson(ctx context.Context, reportChecks []string, reportErrors []string) ([]string, []string, error) {
	var err error = nil

	// Short-circuit exit.

	if len(checkself.EngineConfigurationJson) == 0 {
		return reportChecks, reportErrors, err
	}

	// Check Config path.

	reportChecks = append(reportChecks, fmt.Sprintf("%s=%s", option.EngineConfigurationJson.Envar, checkself.EngineConfigurationJson))

	parsedEngineConfigurationJson := &engineconfigurationjsonparser.EngineConfigurationJsonParserImpl{
		EngineConfigurationJson: checkself.EngineConfigurationJson,
	}

	// Test ConfigPath

	configPath, err := parsedEngineConfigurationJson.GetConfigPath(ctx)
	if err != nil {
		return reportChecks, reportErrors, err
	}

	requiredFiles := []string{
		"cfgVariant.json",
		"defaultGNRCP.config",
	}
	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", configPath, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s=%s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", option.ConfigPath.Envar, checkself.ConfigPath, targetFile))
		}
	}

	return reportChecks, reportErrors, err
}
