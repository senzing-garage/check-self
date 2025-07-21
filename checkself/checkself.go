package checkself

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"

	"github.com/senzing-garage/go-databasing/connector"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-helpers/settingsparser"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicCheckSelf is the basic checker.
type BasicCheckSelf struct {
	ConfigPath                 string
	DatabaseURL                string
	EngineLogLevel             string // IMPROVE:
	ErrorLicenseDaysLeft       string
	ErrorLicenseRecordsPercent string
	GrpcDialOptions            []grpc.DialOption // IMPROVE:
	GrpcURL                    string            // IMPROVE:
	InputURL                   string            // IMPROVE:
	LicenseStringBase64        string            // IMPROVE:
	LogLevel                   string            // IMPROVE:
	ObserverURL                string            // IMPROVE:
	ResourcePath               string
	SenzingDirectory           string // IMPROVE:
	SenzingInstanceName        string
	SenzingVerboseLogging      int64
	Settings                   string
	SupportPath                string
}

type ProductLicenseResponse struct {
	Billing      string `json:"billing"`
	Contract     string `json:"contract"`
	Customer     string `json:"customer"`
	ExpireDate   string `json:"expireDate"`
	IssueDate    string `json:"issueDate"`
	LicenseLevel string `json:"licenseLevel"`
	LicenseType  string `json:"licenseType"`
	RecordLimit  int64  `json:"recordLimit"`
}

const (
	horizontalRuleLength  = 80
	horizontalTitleLength = horizontalRuleLength - 4
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var defaultInstanceName = "check-self"

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The CheckSelf method prints the output of numerous checks.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/
func (checkself *BasicCheckSelf) CheckSelf(ctx context.Context) error {
	var err error

	reportChecks := []string{}
	reportInfo := []string{}
	reportErrors := []string{}

	// List tests.  Order is important.

	testFunctions := checkself.getTestFunctions()

	// Perform checks.

	for _, testFunction := range testFunctions {
		reportChecks, reportInfo, reportErrors, err = testFunction(ctx, reportChecks, reportInfo, reportErrors)
		if err != nil {
			if len(err.Error()) > 0 {
				reportErrors = append(reportErrors, err.Error())
			}

			break
		}
	}

	// Output reports.

	if len(reportInfo) > 0 {
		printTitle("Information")

		for _, message := range reportInfo {
			outputln(message)
		}
	}

	if len(reportChecks) > 0 {
		printTitle("Checks performed")

		for index, message := range reportChecks {
			outputf("%6d. %s\n", index+1, message)
		}
	}

	if len(reportErrors) > 0 {
		printTitle("Errors")

		for index, message := range reportErrors {
			outputf("%6d. %s\n\n", index+1, message)
		}

		err = wraperror.Errorf(errForPackage, "%d errors detected", len(reportErrors))
		outputf("Result: %s\n", err.Error())
	} else {
		printTitle("Result")
		outputf("No errors detected.\n")
	}

	outputf("%s\n\n\n\n\n", strings.Repeat("-", horizontalRuleLength))

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) createSzAbstractFactory(ctx context.Context) (senzing.SzAbstractFactory, error) {
	var (
		err    error
		result senzing.SzAbstractFactory
	)

	result, err = szfactorycreator.CreateCoreAbstractFactory(
		checkself.getInstanceName(ctx),
		checkself.getSettings(ctx),
		checkself.SenzingVerboseLogging,
		senzing.SzInitializeWithDefaultConfiguration,
	)

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (checkself *BasicCheckSelf) createSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	var (
		err    error
		result senzing.SzConfigManager
	)

	szAbstractFactory, err := checkself.createSzAbstractFactory(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not create SzAbstractFactory")
	}

	defer func() {
		err := szAbstractFactory.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()

	result, err = szAbstractFactory.CreateConfigManager(ctx)

	return result, wraperror.Errorf(err, "Could not create SzConfigManager")
}

func (checkself *BasicCheckSelf) createSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	var (
		err    error
		result senzing.SzProduct
	)

	szAbstractFactory, err := checkself.createSzAbstractFactory(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not create SzAbstractFactory")
	}

	defer func() {
		err := szAbstractFactory.Close(ctx)
		if err != nil {
			panic(err)
		}
	}()

	result, err = szAbstractFactory.CreateProduct(ctx)

	return result, wraperror.Errorf(err, "Could not create SzProduct")
}

func (checkself *BasicCheckSelf) getDatabaseURL(ctx context.Context) (string, error) {
	if len(checkself.DatabaseURL) > 0 { // Simple case.
		return checkself.DatabaseURL, nil
	}

	if len(checkself.Settings) == 0 {
		return "", wraperror.Errorf(errForPackage, "neither DatabaseUrl nor settings set")
	}

	// Pull database from Senzing engine configuration json.
	// IMPROVE: This code only returns one database.  Need to handle the multi-database case.

	parsedSettings, err := settingsparser.New(checkself.Settings)
	if err != nil {
		return "", wraperror.Errorf(errForPackage, "unable to parse settings: %s", checkself.Settings)
	}

	databaseUris, err := parsedSettings.GetDatabaseURIs(ctx)
	if err != nil {
		return "", wraperror.Errorf(errForPackage, "unable to extract databases from settings: %s", checkself.Settings)
	}

	if len(databaseUris) == 0 {
		return "", wraperror.Errorf(err, "no databases found in settings: %s", checkself.Settings)
	}

	return databaseUris[0], nil
}

func (checkself *BasicCheckSelf) getDatabaseConnector(ctx context.Context) (driver.Connector, error) {
	var (
		err    error
		result driver.Connector
	)

	databaseURL, err := checkself.getDatabaseURL(ctx)
	if err != nil {
		return result, wraperror.Errorf(
			err,
			"Unable to locate Database URL For more information, visit https://hub.senzing.com/...",
		)
	}

	result, err = connector.NewConnector(ctx, databaseURL)
	if err != nil {
		return result, wraperror.Errorf(
			err,
			"Database URL '%s' is misconfigured. Could not create a database connector. For more information, visit https://hub.senzing.com/...",
			databaseURL,
		)
	}

	return result, wraperror.Errorf(err, wraperror.NoMessage)
}

func (checkself *BasicCheckSelf) getInstanceName(ctx context.Context) string {
	_ = ctx

	result := checkself.SenzingInstanceName
	if len(result) == 0 {
		result = defaultInstanceName
	}

	return result
}

func (checkself *BasicCheckSelf) getSettings(ctx context.Context) string {
	_ = ctx

	var err error

	result := checkself.Settings
	if len(result) == 0 {
		result, err = settings.BuildSimpleSettingsUsingEnvVars()
		if err != nil {
			panic(err.Error())
		}
	}

	return result
}

func (checkself *BasicCheckSelf) getTestFunctions() []func(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {
	return []func(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error){
		checkself.Prolog,
		checkself.ListEnvironmentVariables,
		checkself.ListStructVariables,
		checkself.CheckConfigPath,
		checkself.CheckResourcePath,
		checkself.CheckSupportPath,
		checkself.CheckDatabaseURL,
		checkself.CheckSettings,
		checkself.Break,
		checkself.CheckDatabaseSchema,
		checkself.Break,
		// checkself.CheckSenzingConfiguration,
		// checkself.CheckLicense,
	}
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func statFiles(variableName string, path string, requiredFiles []string) []string {
	reportErrors := []string{}

	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", path, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(
				reportErrors,
				fmt.Sprintf(
					"%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...",
					variableName,
					path,
					targetFile,
				),
			)
		}
	}

	return reportErrors
}

func outputf(format string, message ...any) {
	fmt.Printf(format, message...) //nolint
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}

func printTitle(title string) {
	outputf("\n-- %s %s\n\n", title, strings.Repeat("-", horizontalTitleLength-len(title)))
}
