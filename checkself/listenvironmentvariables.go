package checkself

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) ListEnvironmentVariables(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	_ = ctx

	osEnviron := map[string]string{}
	for _, element := range os.Environ() {
		if strings.HasPrefix(element, "SENZING_TOOLS_") {
			variable := strings.Split(element, "=")
			osEnviron[variable[0]] = variable[1]
		}
	}

	if len(osEnviron) > 0 {
		reportInfo = append(reportInfo, "\nSENZING_TOOLS_* environment variables defined:\n")
		count := 0
		for key, value := range osEnviron {
			count++
			reportInfo = append(reportInfo, fmt.Sprintf("%6d. %s = %s", count, key, value))
		}
		reportInfo = append(reportInfo, "")
	}

	return reportChecks, reportInfo, reportErrors, nil
}
