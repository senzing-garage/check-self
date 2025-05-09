package checkself

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/go-databasing/checker"
	"github.com/senzing-garage/go-helpers/wraperror"
)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) CheckLicense(
	ctx context.Context,
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
) ([]string, []string, []string, error) {
	var err error

	// Prolog.

	reportChecks = append(reportChecks, "Check Senzing license")

	// Connect to the database.

	databaseConnector, err := checkself.getDatabaseConnector(ctx)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Get number of record in DSRC_RECORD.

	recordCount, err := checkself.getRecordCount(ctx, databaseConnector)
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Get license

	szProduct, err := checkself.getSzProduct(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not create szProduct.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	license, err := szProduct.GetLicense(ctx)
	if err != nil {
		reportErrors = append(reportErrors, fmt.Sprintf("Could not get license information.  Error %s", err.Error()))
		return reportChecks, reportInfo, reportErrors, nil
	}

	// Marshal license into structure.

	productLicenseResponse := &ProductLicenseResponse{}
	err = json.Unmarshal([]byte(license), productLicenseResponse)
	if err != nil {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf("Could not parse license information into structure.  Error %s", err.Error()),
		)
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
		reportErrors = append(
			reportErrors,
			fmt.Sprintf("Could not parse expireDate information.  Error %s", err.Error()),
		)
		return reportChecks, reportInfo, reportErrors, nil
	}
	duration := time.Until(licenseExpireDate)
	expireInDays := int(duration.Hours() / 24)

	reportInfo = append(reportInfo, fmt.Sprintf(
		`
License:

- Records used: %d of %d
- Date license expires: %s
- Days until license expires: %d

%s`,
		recordCount,
		productLicenseResponse.RecordLimit,
		productLicenseResponse.ExpireDate,
		expireInDays,
		prettyJSON.String(),
	))

	// Calculate License Days Left error.

	if len(checkself.ErrorLicenseDaysLeft) == 0 {
		checkself.ErrorLicenseDaysLeft = DefaultSenzingToolsLicenseDaysLeft
	}
	errorLicenseDaysLeft, err := strconv.Atoi(checkself.ErrorLicenseDaysLeft)
	if err != nil {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"Could not parse SENZING_TOOLS_LICENSE_DAYS_LEFT information: %s.  Error %s",
				checkself.ErrorLicenseDaysLeft,
				err.Error(),
			),
		)
		return reportChecks, reportInfo, reportErrors, nil
	}
	if expireInDays < errorLicenseDaysLeft {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"License expires in %d days. For more information, visit https://hub.senzing.com/... ",
				expireInDays,
			),
		)
	}

	// Calculate License Records Percent error.

	if len(checkself.ErrorLicenseRecordsPercent) == 0 {
		checkself.ErrorLicenseRecordsPercent = DefaultSenzingToolsLicenseRecordsPercent
	}
	errorLicenseRecordsPercent, err := strconv.Atoi(checkself.ErrorLicenseRecordsPercent)
	if err != nil {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"Could not parse SENZING_TOOLS_LICENSE_RECORDS_PERCENT information: %s.  Error %s",
				checkself.ErrorLicenseRecordsPercent,
				err.Error(),
			),
		)
		return reportChecks, reportInfo, reportErrors, nil
	}
	if (recordCount / productLicenseResponse.RecordLimit) > int64(errorLicenseRecordsPercent) {
		reportErrors = append(
			reportErrors,
			fmt.Sprintf(
				"Records above %d full limit. For more information, visit https://hub.senzing.com/... ",
				errorLicenseRecordsPercent,
			),
		)
	}

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}

func (checkself *BasicCheckSelf) getRecordCount(
	ctx context.Context,
	databaseConnector driver.Connector,
) (int64, error) {
	var (
		err    error
		result int64
	)
	checker := &checker.BasicChecker{
		DatabaseConnector: databaseConnector,
	}
	result, err = checker.RecordCount(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not get count of records.  Error %w", err)
	}

	return result, err
}
