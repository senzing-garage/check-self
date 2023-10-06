package checkself

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/go-databasing/checker"
	"github.com/senzing/go-databasing/connector"
)

func (checkself *CheckSelfImpl) CheckLicense(ctx context.Context, reportChecks []string, reportInfo []string, reportErrors []string) ([]string, []string, []string, error) {

	// Prolog.

	reportChecks = append(reportChecks, "Check Senzing license")

	// Connect to the database.

	databaseConnector, err := connector.NewConnector(ctx, checkself.DatabaseUrl)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("%s = %s is misconfigured. Could not create a databsae connector. For more information, visit https://hub.senzing.com/...  Error: %s", option.DatabaseUrl.Envar, checkself.DatabaseUrl, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Get number of record in DSRC_RECORD.

	checker := &checker.CheckerImpl{
		DatabaseConnector: databaseConnector,
	}
	recordCount, err := checker.RecordCount(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not get count of records.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Get license

	g2Product, err := checkself.getG2product(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not create g2Product.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	license, err := g2Product.License(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not get license information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Marshal license into structure.

	productLicenseResponse := &ProductLicenseResponse{}
	err = json.Unmarshal([]byte(license), productLicenseResponse)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse license information into structure.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Pretty-print JSON.

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(license), "", "\t")
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse license information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	licenseExpireDate, err := time.Parse(time.DateOnly, productLicenseResponse.ExpireDate)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse expireDate information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}
	duration := time.Until(licenseExpireDate)
	expireInDays := int(duration.Hours() / 24)

	reportInfo = append(reportInfo, fmt.Sprintf(`
License:

- Records used: %d of %d
- Date license expires: %s
- Days until license expires: %d

%s`,
		recordCount, productLicenseResponse.RecordLimit, productLicenseResponse.ExpireDate, expireInDays, prettyJSON.String()))

	// Calculate errors.

	errorLicenseDaysLeft, err := strconv.Atoi(checkself.ErrorLicenseDaysLeft)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse SENZING_TOOLS_LICENSE_DAYS_LEFT information: %s.  Error %s", checkself.ErrorLicenseDaysLeft, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	if expireInDays < errorLicenseDaysLeft {
		reportErrors = append(reportErrors, fmt.Sprintf("License expires in %d days. For more information, visit https://hub.senzing.com/... ", expireInDays))
	}

	errorLicenseRecordsPercent, err := strconv.Atoi(checkself.ErrorLicenseRecordsPercent)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not parse SENZING_TOOLS_LICENSE_RECORDS_PERCENT information: %s.  Error %s", checkself.ErrorLicenseRecordsPercent, err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	if (recordCount / productLicenseResponse.RecordLimit) > int64(errorLicenseRecordsPercent) {
		reportErrors = append(reportErrors, fmt.Sprintf("Records above %d full limit. For more information, visit https://hub.senzing.com/... ", errorLicenseRecordsPercent))
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}
