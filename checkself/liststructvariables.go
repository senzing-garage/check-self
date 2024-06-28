package checkself

import (
	"context"
	"fmt"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) ListStructVariables(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	_ = ctx

	structStrings := map[string]string{
		"ConfigPath":          checkself.ConfigPath,
		"DatabaseURL":         checkself.DatabaseURL,
		"Settings":            checkself.Settings,
		"EngineLogLevel":      checkself.EngineLogLevel,
		"GrpcURL":             checkself.GrpcURL,
		"InputURL":            checkself.InputURL,
		"LicenseStringBase64": checkself.LicenseStringBase64,
		"LogLevel":            checkself.LogLevel,
		"ObserverURL":         checkself.ObserverURL,
		"ResourcePath":        checkself.ResourcePath,
		"SenzingDirectory":    checkself.SenzingDirectory,
		"SupportPath":         checkself.SupportPath,
	}

	count := 0
	reportInfo = append(reportInfo, "\nCommand line variables:\n")

	for key, value := range structStrings {
		if len(value) > 0 {
			count++
			reportInfo = append(reportInfo, fmt.Sprintf("%6d. %s = %s", count, key, value))
		}
	}
	reportInfo = append(reportInfo, "")

	return reportChecks, reportInfo, reportErrors, nil
}
