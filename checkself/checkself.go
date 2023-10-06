package checkself

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-sdk-abstract-factory/factory"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// CheckSelfImpl is the basic checker.
type CheckSelfImpl struct {
	// g2configSingleton       g2api.G2config
	// g2configSyncOnce        sync.Once
	ConfigPath                 string
	DatabaseUrl                string
	EngineConfigurationJson    string
	EngineLogLevel             string // TODO:
	ErrorLicenseDaysLeft       string
	ErrorLicenseRecordsPercent string
	g2configmgrSingleton       g2api.G2configmgr
	g2configmgrSyncOnce        sync.Once
	g2factorySingleton         factory.SdkAbstractFactory
	g2factorySyncOnce          sync.Once
	g2productSingleton         g2api.G2product
	g2productSyncOnce          sync.Once
	GrpcDialOptions            []grpc.DialOption // TODO:
	GrpcUrl                    string            // TODO:
	InputUrl                   string            // TODO:
	LicenseStringBase64        string            // TODO:
	LogLevel                   string            // TODO:
	ObserverUrl                string            // TODO:
	ResourcePath               string
	SenzingDirectory           string // TODO:
	SenzingModuleName          string
	SenzingVerboseLogging      int
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

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var defaultModuleName string = "check-self"

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

func (checkself *CheckSelfImpl) getEngineConfigurationJson(ctx context.Context) string {
	var err error = nil
	result := checkself.EngineConfigurationJson
	if len(result) == 0 {
		result, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingEnvVars()
		if err != nil {
			panic(err.Error())
		}
	}
	return result
}

func (checkself *CheckSelfImpl) getModuleName(ctx context.Context) string {
	result := checkself.SenzingModuleName
	if len(result) == 0 {
		result = defaultModuleName
	}
	return result
}

// Create a G2Configmgr singleton and return it.
func (checkself *CheckSelfImpl) getG2configmgr(ctx context.Context) (g2api.G2configmgr, error) {
	var err error = nil
	checkself.g2configmgrSyncOnce.Do(func() {
		checkself.g2configmgrSingleton, err = checkself.getG2Factory(ctx).GetG2configmgr(ctx)
		if err != nil {
			return
		}
		if checkself.g2configmgrSingleton.GetSdkId(ctx) == "base" {
			err = checkself.g2configmgrSingleton.Init(ctx, checkself.getModuleName(ctx), checkself.getEngineConfigurationJson(ctx), checkself.SenzingVerboseLogging)
		}
	})
	return checkself.g2configmgrSingleton, err
}

func (checkself *CheckSelfImpl) getG2Factory(ctx context.Context) factory.SdkAbstractFactory {
	checkself.g2factorySyncOnce.Do(func() {
		checkself.g2factorySingleton = &factory.SdkAbstractFactoryImpl{}
	})
	return checkself.g2factorySingleton
}

// Create a G2Configmgr singleton and return it.
func (checkself *CheckSelfImpl) getG2product(ctx context.Context) (g2api.G2product, error) {
	var err error = nil
	checkself.g2productSyncOnce.Do(func() {
		checkself.g2productSingleton, err = checkself.getG2Factory(ctx).GetG2product(ctx)
		if err != nil {
			return
		}
		if checkself.g2configmgrSingleton.GetSdkId(ctx) == "base" {
			err = checkself.g2productSingleton.Init(ctx, checkself.getModuleName(ctx), checkself.getEngineConfigurationJson(ctx), checkself.SenzingVerboseLogging)
		}
	})
	return checkself.g2productSingleton, err
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
	fmt.Printf("%s\n\n", strings.Repeat("-", 80))

	return err
}
