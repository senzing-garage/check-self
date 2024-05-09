package checkself

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/senzing-garage/go-helpers/engineconfigurationjson"
	"github.com/senzing-garage/go-helpers/engineconfigurationjsonparser"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/sz"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// CheckSelfImpl is the basic checker.
type CheckSelfImpl struct {
	ConfigPath                 string
	DatabaseUrl                string
	EngineLogLevel             string // TODO:
	ErrorLicenseDaysLeft       string
	ErrorLicenseRecordsPercent string
	GrpcDialOptions            []grpc.DialOption // TODO:
	GrpcUrl                    string            // TODO:
	InputUrl                   string            // TODO:
	LicenseStringBase64        string            // TODO:
	LogLevel                   string            // TODO:
	ObserverUrl                string            // TODO:
	ResourcePath               string
	SenzingDirectory           string // TODO:
	SenzingInstanceName        string
	SenzingVerboseLogging      int64
	Settings                   string
	SupportPath                string
	szConfigManagerSingleton   sz.SzConfigManager
	szConfigManagerSyncOnce    sync.Once
	szFactorySingleton         sz.SzAbstractFactory
	szFactorySyncOnce          sync.Once
	szProductSingleton         sz.SzProduct
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

var defaultInstanceName string = "check-self"

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func statFiles(variableName string, path string, requiredFiles []string) []string {
	reportErrors := []string{}
	for _, requiredFile := range requiredFiles {
		targetFile := fmt.Sprintf("%s/%s", path, requiredFile)
		if _, err := os.Stat(targetFile); err != nil {
			reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not find %s. For more information, visit https://hub.senzing.com/...", variableName, path, targetFile))
		}
	}
	return reportErrors
}

func printTitle(title string) {
	fmt.Printf("\n-- %s %s\n\n", title, strings.Repeat("-", 76-len(title)))
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (checkself *CheckSelfImpl) getDatabaseUrl(ctx context.Context) (string, error) {

	// Simple case.

	if len(checkself.DatabaseUrl) > 0 {
		return checkself.DatabaseUrl, nil
	}

	if len(checkself.Settings) == 0 {
		return "", fmt.Errorf("neither DatabaseUrl nor EngineConfigurationJson set")
	}

	// Pull database from Senzing engine configuration json.
	// TODO: This code only returns one database.  Need to handle the multi-database case.

	parsedEngineConfigurationJson, err := engineconfigurationjsonparser.New(checkself.Settings)
	if err != nil {
		return "", fmt.Errorf("unable to parse EngineConfigurationJson: %s", checkself.Settings)
	}

	databaseUrls, err := parsedEngineConfigurationJson.GetDatabaseUrls(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to extract databases from EngineConfigurationJson: %s", checkself.Settings)
	}
	if len(databaseUrls) == 0 {
		return "", fmt.Errorf("no databases found in EngineConfigurationJson: %s", checkself.Settings)
	}
	return databaseUrls[0], nil
}

func (checkself *CheckSelfImpl) getSettings(ctx context.Context) string {
	_ = ctx
	var err error = nil
	result := checkself.Settings
	if len(result) == 0 {
		result, err = engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			panic(err.Error())
		}
	}
	return result
}

func (checkself *CheckSelfImpl) getInstanceName(ctx context.Context) string {
	_ = ctx
	result := checkself.SenzingInstanceName
	if len(result) == 0 {
		result = defaultInstanceName
	}
	return result
}

// Create a SzConfigManager singleton and return it.
func (checkself *CheckSelfImpl) getSzConfigManager(ctx context.Context) (sz.SzConfigManager, error) {
	var err error = nil
	checkself.szConfigManagerSyncOnce.Do(func() {
		checkself.szConfigManagerSingleton, err = checkself.getSzFactory(ctx).CreateSzConfigManager(ctx)
		if err != nil {
			return
		}
	})
	return checkself.szConfigManagerSingleton, err
}

func (checkself *CheckSelfImpl) getSzFactory(ctx context.Context) sz.SzAbstractFactory {
	var err error = nil
	checkself.szFactorySyncOnce.Do(func() {
		checkself.szFactorySingleton, err = szfactorycreator.CreateCoreAbstractFactory(checkself.getInstanceName(ctx), checkself.getSettings(ctx), checkself.SenzingVerboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	})
	if err != nil {
		return nil
	}
	return checkself.szFactorySingleton
}

// Create a SzProduct singleton and return it.
func (checkself *CheckSelfImpl) getSzProduct(ctx context.Context) (sz.SzProduct, error) {
	var err error = nil
	checkself.szProductSyncOnce.Do(func() {
		checkself.szProductSingleton, err = checkself.getSzFactory(ctx).CreateSzProduct(ctx)
		if err != nil {
			return
		}
	})
	return checkself.szProductSingleton, err
}

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
func (checkself *CheckSelfImpl) CheckSelf(ctx context.Context) error {
	var err error = nil

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
		checkself.CheckDatabaseUrl,
		checkself.CheckEngineConfigurationJson,
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
		err = fmt.Errorf("%d errors detected", len(reportErrors))
		fmt.Printf("Result: %s\n", err.Error())
	} else {
		printTitle("Result")
		fmt.Printf("No errors detected.\n")
	}
	fmt.Printf("%s\n\n\n\n\n", strings.Repeat("-", 80))

	return err
}
