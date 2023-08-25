package checkself

import (
	"context"
	"fmt"
	"strings"

	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-common/engineconfigurationjsonparser"
	"github.com/senzing/go-common/g2engineconfigurationjson"
)

// ----------------------------------------------------------------------------
// Helper methods
// ----------------------------------------------------------------------------

func (checkself *CheckSelfImpl) buildAndCheckEngineConfigurationJson(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	engineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
	if err != nil {
		return reportChecks, reportInfo, reportErrors, err
	}
	reportInfo = append(reportInfo, fmt.Sprintf("\nexport SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='%s'\n", engineConfigurationJson))
	return checkself.checkEngineConfigurationJson(ctx, engineConfigurationJson, reportChecks, reportInfo, reportErrors)
}

func (checkself *CheckSelfImpl) checkEngineConfigurationJson(ctx context.Context, engineConfigurationJson string, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	var err error = nil

	parsedEngineConfigurationJson := &engineconfigurationjsonparser.EngineConfigurationJsonParserImpl{
		EngineConfigurationJson: engineConfigurationJson,
	}

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH.

	configVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH"
	configValue, err := parsedEngineConfigurationJson.GetConfigPath(ctx)
	if err != nil {
		return reportChecks, reportInfo, reportErrors, err
	}
	errorList := statFiles(configVariable, configValue, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH.

	resourceVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH"
	resourceValue, err := parsedEngineConfigurationJson.GetResourcePath(ctx)
	if err != nil {
		return reportChecks, reportInfo, reportErrors, err
	}
	errorList = statFiles(resourceVariable, resourceValue, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH.

	supportVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH"
	supportValue, err := parsedEngineConfigurationJson.GetSupportPath(ctx)
	if err != nil {
		return reportChecks, reportInfo, reportErrors, err
	}
	errorList = statFiles(supportVariable, supportValue, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION.

	connectionVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION"
	connectionValues, err := parsedEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		return reportChecks, reportInfo, reportErrors, err
	}
	for _, connectionValue := range connectionValues {
		errorList = checkDatabaseUrl(ctx, connectionVariable, connectionValue)
		reportErrors = append(reportErrors, errorList...)
	}

	return reportChecks, reportInfo, reportErrors, err

}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *CheckSelfImpl) CheckEngineConfigurationJson(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	if len(checkself.EngineConfigurationJson) == 0 {
		return checkself.buildAndCheckEngineConfigurationJson(ctx, reportChecks, reportInfo, reportErrors)
	}

	// Verify that JSON string is syntactically correct.

	parsedEngineConfigurationJson, err := engineconfigurationjsonparser.New(checkself.EngineConfigurationJson)
	if err != nil {
		normalizedValue := strings.ReplaceAll(strings.ReplaceAll(checkself.EngineConfigurationJson, "\n", " "), "  ", "")
		reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineConfigurationJson.Envar, normalizedValue))
		reportErrors = append(reportErrors, fmt.Sprintf("%s - %s", option.EngineConfigurationJson.Envar, err.Error()))
		return reportChecks, reportInfo, reportErrors, err
	}

	// Report what is being checked.

	redactedJson, err := parsedEngineConfigurationJson.RedactedJson(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, err
	}
	reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineConfigurationJson.Envar, redactedJson))

	// Perform check.

	return checkself.checkEngineConfigurationJson(ctx, checkself.EngineConfigurationJson, reportChecks, reportInfo, reportErrors)
}
