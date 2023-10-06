package checkself

import (
	"bytes"
	"context"
	"encoding/json"
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
		reportErrors = append(reportErrors, fmt.Sprintf("Could not build engine configuration json. %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(engineConfigurationJson), "", "\t")
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse license information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	reportInfo = append(reportInfo, fmt.Sprintf("\nEffective engine configuration:\n\nexport SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='%s'\n", prettyJSON.String()))
	return checkself.checkEngineConfigurationJson(ctx, engineConfigurationJson, reportChecks, reportInfo, reportErrors)
}

func (checkself *CheckSelfImpl) checkEngineConfigurationJson(ctx context.Context, engineConfigurationJson string, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	parsedEngineConfigurationJson := &engineconfigurationjsonparser.EngineConfigurationJsonParserImpl{
		EngineConfigurationJson: engineConfigurationJson,
	}

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH.

	configVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH"
	configValue, err := parsedEngineConfigurationJson.GetConfigPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", configVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList := statFiles(configVariable, configValue, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH.

	resourceVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH"
	resourceValue, err := parsedEngineConfigurationJson.GetResourcePath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", resourceVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList = statFiles(resourceVariable, resourceValue, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH.

	supportVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH"
	supportValue, err := parsedEngineConfigurationJson.GetSupportPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", supportVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList = statFiles(supportVariable, supportValue, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION.

	connectionVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION"
	connectionValues, err := parsedEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", connectionVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	for _, connectionValue := range connectionValues {
		errorList = checkDatabaseUrl(ctx, connectionVariable, connectionValue)
		reportErrors = append(reportErrors, errorList...)
	}

	return reportChecks, reportInfo, reportErrors, nil

}

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *CheckSelfImpl) CheckEngineConfigurationJson(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.EngineConfigurationJson) == 0 {
		return checkself.buildAndCheckEngineConfigurationJson(ctx, reportChecks, reportInfo, reportErrors)
	}

	// Verify that JSON string is syntactically correct.

	parsedEngineConfigurationJson, err := engineconfigurationjsonparser.New(checkself.EngineConfigurationJson)
	if err != nil {
		normalizedValue := strings.ReplaceAll(strings.ReplaceAll(checkself.EngineConfigurationJson, "\n", " "), "  ", "")
		reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineConfigurationJson.Envar, normalizedValue))
		reportErrors = append(reportErrors, fmt.Sprintf("%s - %s", option.EngineConfigurationJson.Envar, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	databaseUrls, err := parsedEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}
	for _, databaseUrl := range databaseUrls {
		errorList := checkDatabaseUrl(ctx, option.EngineConfigurationJson.Envar, databaseUrl)
		reportErrors = append(reportErrors, errorList...)
	}

	// Report what is being checked.

	redactedJson, err := parsedEngineConfigurationJson.RedactedJson(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}
	reportChecks = append(reportChecks, fmt.Sprintf("Check engine configuration: %s = %s", option.EngineConfigurationJson.Envar, redactedJson))

	// Perform check.

	return checkself.checkEngineConfigurationJson(ctx, checkself.EngineConfigurationJson, reportChecks, reportInfo, reportErrors)
}
