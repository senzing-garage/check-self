package checkself

import (
	"context"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckSenzingConfiguration(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	reportChecks = append(reportChecks, "Check Senzing configuration")

	// Create Senzing objects.

	szConfigManager, err := checkself.getSzConfigManager(ctx)
	if err != nil {
		reportErrors = append(reportErrors, "Could not create szConfigManager.  Error: "+err.Error())

		return reportChecks, reportInfo, reportErrors, nil //nolint
	}

	// Determine if Configuration exists.

	configID, err := szConfigManager.GetDefaultConfigID(ctx)
	if err != nil {
		reportErrors = append(
			reportErrors,
			"Could not get Senzing default configuration ID.  Error "+err.Error(),
		)

		return reportChecks, reportInfo, reportErrors, nil //nolint
	}

	if configID == 0 {
		reportErrors = append(
			reportErrors,
			"Senzing configuration doesn't exist. For more information, visit https://hub.senzing.com/...",
		)
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
