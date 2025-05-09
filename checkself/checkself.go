package checkself

import (
	"context"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"sync"

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
	EngineLogLevel             string // TODO:
	ErrorLicenseDaysLeft       string
	ErrorLicenseRecordsPercent string
	GrpcDialOptions            []grpc.DialOption // TODO:
	GrpcURL                    string            // TODO:
	InputURL                   string            // TODO:
	LicenseStringBase64        string            // TODO:
	LogLevel                   string            // TODO:
	ObserverURL                string            // TODO:
	ResourcePath               string
	SenzingDirectory           string // TODO:
	SenzingInstanceName        string
	SenzingVerboseLogging      int64
	Settings                   string
	SupportPath                string
	szConfigManagerSingleton   senzing.SzConfigManager
	szConfigManagerSyncOnce    sync.Once
	szFactorySingleton         senzing.SzAbstractFactory
	szFactorySyncOnce          sync.Once
	szProductSingleton         senzing.SzProduct
	szProductSyncOnce          sync.Once
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

	testFunctions := []func(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error){
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
		checkself.CheckSenzingConfiguration,
		checkself.CheckLicense,
	}

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
			fmt.Println(message)
		}
	}

	if len(reportChecks) > 0 {
		printTitle("Checks performed")
		for index, message := range reportChecks {
			fmt.Printf("%6d. %s\n", index+1, message)
		}
	}

	if len(reportErrors) > 0 {
		printTitle("Errors")
		for index, message := range reportErrors {
			fmt.Printf("%6d. %s\n\n", index+1, message)
		}
		err = wraperror.Errorf(errForPackage, "%d errors detected", len(reportErrors))
		fmt.Printf("Result: %s\n", err.Error())
	} else {
		printTitle("Result")
		fmt.Printf("No errors detected.\n")
	}
	fmt.Printf("%s\n\n\n\n\n", strings.Repeat("-", 80))

	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) getDatabaseURL(ctx context.Context) (string, error) {

	// Simple case.

	if len(checkself.DatabaseURL) > 0 {
		return checkself.DatabaseURL, nil
	}

	if len(checkself.Settings) == 0 {
		return "", wraperror.Errorf(errForPackage, "neither DatabaseUrl nor settings set")
	}

	// Pull database from Senzing engine configuration json.
	// TODO: This code only returns one database.  Need to handle the multi-database case.

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
			"Unable to locate Database URL For more information, visit https://hub.senzing.com/...  Error: %w",
			err,
		)
	}

	result, err = connector.NewConnector(ctx, databaseURL)
	if err != nil {
		return result, wraperror.Errorf(
			err,
			"Database URL '%s' is misconfigured. Could not create a database connector. For more information, visit https://hub.senzing.com/...  Error: %w",
			databaseURL,
			err,
		)

	}

	return result, err
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

func (checkself *BasicCheckSelf) getInstanceName(ctx context.Context) string {
	_ = ctx
	result := checkself.SenzingInstanceName
	if len(result) == 0 {
		result = defaultInstanceName
	}
	return result
}

// Create a SzConfigManager singleton and return it.
func (checkself *BasicCheckSelf) getSzConfigManager(ctx context.Context) (senzing.SzConfigManager, error) {
	var err error
	checkself.szConfigManagerSyncOnce.Do(func() {
		checkself.szConfigManagerSingleton, err = checkself.getSzFactory(ctx).CreateConfigManager(ctx)
	})
	return checkself.szConfigManagerSingleton, err
}

func (checkself *BasicCheckSelf) getSzFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error
	checkself.szFactorySyncOnce.Do(func() {
		checkself.szFactorySingleton, err = szfactorycreator.CreateCoreAbstractFactory(
			checkself.getInstanceName(ctx),
			checkself.getSettings(ctx),
			checkself.SenzingVerboseLogging,
			senzing.SzInitializeWithDefaultConfiguration,
		)
	})
	if err != nil {
		panic(err.Error())
	}
	return checkself.szFactorySingleton
}

// Create a SzProduct singleton and return it.
func (checkself *BasicCheckSelf) getSzProduct(ctx context.Context) (senzing.SzProduct, error) {
	var err error
	checkself.szProductSyncOnce.Do(func() {
		checkself.szProductSingleton, err = checkself.getSzFactory(ctx).CreateProduct(ctx)
	})
	return checkself.szProductSingleton, err
}

// ----------------------------------------------------------------------------
// Internal functions
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

func printTitle(title string) {
	fmt.Printf("\n-- %s %s\n\n", title, strings.Repeat("-", 76-len(title)))
}
