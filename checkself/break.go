package checkself

import (
	"context"
	"fmt"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) Break(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	_ = ctx
	if len(reportErrors) > 0 {
		return reportChecks, reportInfo, reportErrors, fmt.Errorf("")
	}
	return reportChecks, reportInfo, reportErrors, nil
}
