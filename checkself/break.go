package checkself

import (
	"context"
	"fmt"
)

func (checkself *CheckSelfImpl) Break(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	if len(reportErrors) > 0 {
		return reportChecks, reportInfo, reportErrors, fmt.Errorf("")
	}
	return reportChecks, reportInfo, reportErrors, nil
}
