package checkself

import (
	"context"
	"fmt"
	"strings"

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

	normalizedValue := strings.ReplaceAll(strings.ReplaceAll(checkself.EngineConfigurationJson, "\n", " "), "  ", "")
	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineConfigurationJson.Envar, normalizedValue))

	parsedEngineConfigurationJson := &engineconfigurationjsonparser.EngineConfigurationJsonParserImpl{
		EngineConfigurationJson: checkself.EngineConfigurationJson,
	}

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH.

	configVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH"
	configValue, err := parsedEngineConfigurationJson.GetConfigPath(ctx)
	if err != nil {
		return reportChecks, reportErrors, err
	}
	errorList := statFiles(configVariable, configValue, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH.

	resourceVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH"
	resourceValue, err := parsedEngineConfigurationJson.GetResourcePath(ctx)
	if err != nil {
		return reportChecks, reportErrors, err
	}
	errorList = statFiles(resourceVariable, resourceValue, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH.

	supportVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH"
	supportValue, err := parsedEngineConfigurationJson.GetSupportPath(ctx)
	if err != nil {
		return reportChecks, reportErrors, err
	}
	errorList = statFiles(supportVariable, supportValue, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION.

	connectionVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION"
	connectionValues, err := parsedEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		return reportChecks, reportErrors, err
	}
	for _, connectionValue := range connectionValues {
		errorList = checkDatabaseUrl(ctx, connectionVariable, connectionValue)
		reportErrors = append(reportErrors, errorList...)
	}

	return reportChecks, reportErrors, err
}
