package checkself

import (
	"context"

	"github.com/senzing-garage/go-helpers/wraperror"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) Break(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	_ = ctx
	if len(reportErrors) > 0 {
		return reportChecks, reportInfo, reportErrors, wraperror.Errorf(errForPackage, "")
	}
	return reportChecks, reportInfo, reportErrors, nil
}
