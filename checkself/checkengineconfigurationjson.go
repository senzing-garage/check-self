package checkself

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/settingsparser"
)

// ----------------------------------------------------------------------------
// Helper methods
// ----------------------------------------------------------------------------

func (checkself *CheckSelfImpl) buildAndChecksettings(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not build engine configuration json. %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(settings), "", "\t")
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse license information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	reportInfo = append(reportInfo, fmt.Sprintf("\nEffective engine configuration:\n\nexport SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='%s'\n", prettyJSON.String()))
	return checkself.checksettings(ctx, settings, reportChecks, reportInfo, reportErrors)
}

func (checkself *CheckSelfImpl) checksettings(ctx context.Context, settings string, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	parsedsettings := &settingsparser.BasicSettingsParser{
		Settings: settings,
	}

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH.

	configVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH"
	configValue, err := parsedsettings.GetConfigPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", configVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList := statFiles(configVariable, configValue, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH.

	resourceVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH"
	resourceValue, err := parsedsettings.GetResourcePath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", resourceVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList = statFiles(resourceVariable, resourceValue, RequiredResourceFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH.

	supportVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH"
	supportValue, err := parsedsettings.GetSupportPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", supportVariable, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList = statFiles(supportVariable, supportValue, RequiredSupportFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION.

	connectionVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION"
	connectionValues, err := parsedsettings.GetDatabaseURLs(ctx)
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

func (checkself *CheckSelfImpl) Checksettings(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Short-circuit exit.

	if len(checkself.Settings) == 0 {
		return checkself.buildAndChecksettings(ctx, reportChecks, reportInfo, reportErrors)
	}

	// Verify that JSON string is syntactically correct.

	parsedsettings, err := settingsparser.New(checkself.Settings)
	if err != nil {
		normalizedValue := strings.ReplaceAll(strings.ReplaceAll(checkself.Settings, "\n", " "), "  ", "")
		reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineConfigurationJSON.Envar, normalizedValue))
		reportErrors = append(reportErrors, fmt.Sprintf("%s - %s", option.EngineConfigurationJSON.Envar, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	databaseUrls, err := parsedsettings.GetDatabaseURLs(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}
	for _, databaseUrl := range databaseUrls {
		errorList := checkDatabaseUrl(ctx, option.EngineConfigurationJSON.Envar, databaseUrl)
		reportErrors = append(reportErrors, errorList...)
	}

	// Report what is being checked.

	redactedJson, err := parsedsettings.RedactedJSON(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}
	reportChecks = append(reportChecks, fmt.Sprintf("Check engine configuration: %s = %s", option.EngineConfigurationJSON.Envar, redactedJson))

	// Perform check.

	return checkself.checksettings(ctx, checkself.Settings, reportChecks, reportInfo, reportErrors)
}
