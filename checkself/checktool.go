package checkself

import (
	"context"
	"fmt"
	"time"
)

func (checkself *CheckSelfImpl) CheckTool(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	var err error = nil
	reportInfo = append(reportInfo, fmt.Sprintf("Date: %s ", time.Now().UTC().Format(time.RFC3339)))
	reportInfo = append(reportInfo, fmt.Sprintf("Version: %s-%s ", githubVersion, githubIteration))
	return reportChecks, reportInfo, reportErrors, err
}
