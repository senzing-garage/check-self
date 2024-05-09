package checkself

import (
	"context"
	"fmt"
)

func (checkself *CheckSelfImpl) CheckSenzingConfiguration(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Prolog.

	reportChecks = append(reportChecks, "Check Senzing configuration")

	// Create Senzing objects.

	szConfigManager, err := checkself.getSzConfigManager(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not create szConfigManager.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Determine if Configuration exists.

	configID, err := szConfigManager.GetDefaultConfigId(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not get Senzing default configuration ID.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	if configID == 0 {
		reportErrors = append(reportErrors, "Senzing configuration doesn't exist. For more information, visit https://hub.senzing.com/...")
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
