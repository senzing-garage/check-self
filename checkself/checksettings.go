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
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckSettings(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	if len(checkself.Settings) == 0 { // Short-circuit exit.
		return checkself.buildAndCheckSettings(ctx, reportChecks, reportInfo, reportErrors)
	}

	// Verify that JSON string is syntactically correct.

	parsedSettings, err := settingsparser.New(checkself.Settings)
	if err != nil {
		normalizedValue := strings.ReplaceAll(strings.ReplaceAll(checkself.Settings, "\n", " "), "  ", "")
		reportChecks = append(reportChecks, fmt.Sprintf("%s = %s", option.EngineSettings.Envar, normalizedValue))
		reportErrors = append(reportErrors, fmt.Sprintf("%s - %s", option.EngineSettings.Envar, err.Error()))

		return reportChecks, reportInfo, reportErrors, nil
	}

	databaseURIs, err := parsedSettings.GetDatabaseURIs(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())

		return reportChecks, reportInfo, reportErrors, nil
	}
	for _, databaseURI := range databaseURIs {
		errorList := checkDatabaseURL(ctx, option.EngineSettings.Envar, databaseURI)
		reportErrors = append(reportErrors, errorList...)
	}

	// Report what is being checked.

	redactedJSON, err := parsedSettings.RedactedJSON(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())

		return reportChecks, reportInfo, reportErrors, nil
	}
	reportChecks = append(
		reportChecks,
		fmt.Sprintf("Check engine configuration: %s = %s", option.EngineSettings.Envar, redactedJSON),
	)

	// Perform check.

	return checkself.checkSettings(ctx, checkself.Settings, reportChecks, reportInfo, reportErrors)
}

// ----------------------------------------------------------------------------
// Helper methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) buildAndCheckSettings(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	if err != nil {
		reportErrors = append(reportErrors, "Could not build engine configuration json. "+err.Error())

		return reportChecks, reportInfo, reportErrors, nil //nolint
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(settings), "", "\t")
	if err != nil {
		reportErrors = append(reportErrors, "Could not parse license information.  Error: "+err.Error())

		return reportChecks, reportInfo, reportErrors, nil //nolint
	}

	reportInfo = append(
		reportInfo,
		fmt.Sprintf(
			"\nEffective engine configuration:\n\nexport SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='%s'\n",
			prettyJSON.String(),
		),
	)

	return checkself.checkSettings(ctx, settings, reportChecks, reportInfo, reportErrors)
}

func (checkself *BasicCheckSelf) checkSettings(
	ctx context.Context,
	settings string,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	parsedSettings := &settingsparser.BasicSettingsParser{
		Settings: settings,
	}

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH.

	configVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.CONFIGPATH"
	configValue, err := parsedSettings.GetConfigPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", configVariable, err.Error()))

		return reportChecks, reportInfo, reportErrors, nil
	}
	errorList := statFiles(configVariable, configValue, RequiredConfigFiles)
	reportErrors = append(reportErrors, errorList...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH.

	resourceVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.RESOURCEPATH"
	resourceValue, err := parsedSettings.GetResourcePath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", resourceVariable, err.Error()))

		return reportChecks, reportInfo, reportErrors, nil
	}
	reportErrors = append(reportErrors, statFiles(resourceVariable, resourceValue, RequiredResourceFiles)...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH.

	supportVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.PIPELINE.SUPPORTPATH"
	supportValue, err := parsedSettings.GetSupportPath(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse %s. Error: %s", supportVariable, err.Error()))

		return reportChecks, reportInfo, reportErrors, nil
	}
	reportErrors = append(reportErrors, statFiles(supportVariable, supportValue, RequiredSupportFiles)...)

	// Test SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION.

	connectionVariable := "SENZING_TOOLS_ENGINE_CONFIGURATION_JSON.SQL.CONNECTION"
	databaseURIs, err := parsedSettings.GetDatabaseURIs(ctx)
	if err != nil {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf("Could not parse %s. Error: %s", connectionVariable, err.Error()),
		)

		return reportChecks, reportInfo, reportErrors, nil
	}
	for _, databaseURI := range databaseURIs {
		reportErrors = append(reportErrors, checkDatabaseURL(ctx, connectionVariable, databaseURI)...)
	}

	return reportChecks, reportInfo, reportErrors, nil
}
