package checkself

import (
	"context"
	"fmt"
)

func (checkself *CheckSelfImpl) ListStructVariables(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	structStrings := map[string]string{
		"ConfigPath":          checkself.ConfigPath,
		"DatabaseUrl":         checkself.DatabaseUrl,
		"Settings":            checkself.Settings,
		"EngineLogLevel":      checkself.EngineLogLevel,
		"GrpcUrl":             checkself.GrpcUrl,
		"InputUrl":            checkself.InputUrl,
		"LicenseStringBase64": checkself.LicenseStringBase64,
		"LogLevel":            checkself.LogLevel,
		"ObserverUrl":         checkself.ObserverUrl,
		"ResourcePath":        checkself.ResourcePath,
		"SenzingDirectory":    checkself.SenzingDirectory,
		"SupportPath":         checkself.SupportPath,
	}

	count := 0
	reportInfo = append(reportInfo, "\nCommand line variables:\n")

	for key, value := range structStrings {
		if len(value) > 0 {
			count += 1
			reportInfo = append(reportInfo, fmt.Sprintf("%6d. %s = %s", count, key, value))
		}
	}
	reportInfo = append(reportInfo, "")

	return reportChecks, reportInfo, reportErrors, nil
}
