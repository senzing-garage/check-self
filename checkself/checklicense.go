package checkself

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/senzing-garage/go-databasing/checker"
	"github.com/senzing-garage/go-helpers/wraperror"
)

const (
	hoursPerDay = 24
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

	reportChecks = append(reportChecks, "Check Senzing license")

	recordCount, err := checkself.getRecordCount(ctx)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	license, err := checkself.getLicense(ctx)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	productLicenseResponse, err := getProductLicenseResponse(license)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	prettyJSON, err := getPrettyJSON(license)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	expireInDays, err := getExpireInDays(productLicenseResponse)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	expiryErrors, err := checkself.checkExpiry(expireInDays)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	recordPercentErrors, err := checkself.checkRecordPercent(recordCount, productLicenseResponse)
	if err != nil {
		return returnValues(reportChecks, reportInfo, reportErrors, err)
	}

	reportInfo = append(reportInfo,
		buildReportInfo(
			recordCount,
			productLicenseResponse,
			expireInDays,
			prettyJSON.String(),
		)...)

	reportErrors = append(reportErrors, expiryErrors...)
	reportErrors = append(reportErrors, recordPercentErrors...)

	// Epilog.

	return reportChecks, reportInfo, reportErrors, nil
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (checkself *BasicCheckSelf) checkExpiry(expireInDays int) ([]string, error) {
	var result []string

	if len(checkself.ErrorLicenseDaysLeft) == 0 {
		checkself.ErrorLicenseDaysLeft = DefaultSenzingToolsLicenseDaysLeft
	}
	errorLicenseDaysLeft, err := strconv.Atoi(checkself.ErrorLicenseDaysLeft)
	if err != nil {
		return result, wraperror.Errorf(
			err,
			"Could not parse SENZING_TOOLS_LICENSE_DAYS_LEFT information: %s.  error: %w",
			checkself.ErrorLicenseDaysLeft,
			err,
		)
	}
	if expireInDays < errorLicenseDaysLeft {
		result = append(
			result,
			fmt.Sprintf(
				"License expires in %d days. For more information, visit https://hub.senzing.com/... ",
				expireInDays,
			),
		)
	}

	return result, err
}

func (checkself *BasicCheckSelf) getLicense(ctx context.Context) (string, error) {
	var (
		err    error
		result string
	)

	szProduct, err := checkself.getSzProduct(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not create szProduct.  error: %w", err)
	}

	result, err = szProduct.GetLicense(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not get license information.  error: %w", err)
	}

	return result, err
}

func (checkself *BasicCheckSelf) checkRecordPercent(
	recordCount int64,
	productLicenseResponse *ProductLicenseResponse,
) ([]string, error) {
	var result []string

	if len(checkself.ErrorLicenseRecordsPercent) == 0 {
		checkself.ErrorLicenseRecordsPercent = DefaultSenzingToolsLicenseRecordsPercent
	}
	errorLicenseRecordsPercent, err := strconv.Atoi(checkself.ErrorLicenseRecordsPercent)
	if err != nil {
		return result, wraperror.Errorf(
			err,
			"Could not parse SENZING_TOOLS_LICENSE_RECORDS_PERCENT information: %s.  error: %w",
			checkself.ErrorLicenseRecordsPercent,
			err,
		)
	}

	if (recordCount / productLicenseResponse.RecordLimit) > int64(errorLicenseRecordsPercent) {
		result = append(
			result,
			fmt.Sprintf(
				"Records above %d full limit. For more information, visit https://hub.senzing.com/... ",
				errorLicenseRecordsPercent,
			),
		)
	}

	return result, err
}

func (checkself *BasicCheckSelf) getRecordCount(
	ctx context.Context,
) (int64, error) {
	var (
		err    error
		result int64
	)

	databaseConnector, err := checkself.getDatabaseConnector(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not connect to database.  Error %w", err)
	}

	checker := &checker.BasicChecker{
		DatabaseConnector: databaseConnector,
	}
	result, err = checker.RecordCount(ctx)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not get count of records.  Error %w", err)
	}

	return result, err
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func buildReportInfo(
	recordCount int64,
	productLicenseResponse *ProductLicenseResponse,
	expireInDays int,
	prettyJSON string,
) []string {
	result := []string{
		fmt.Sprintf(`
	License:

	- Records used: %d of %d
	- Date license expires: %s
	- Days until license expires: %d

	%s`,
			recordCount,
			productLicenseResponse.RecordLimit,
			productLicenseResponse.ExpireDate,
			expireInDays,
			prettyJSON,
		)}

	return result
}

func getExpireInDays(productLicenseResponse *ProductLicenseResponse) (int, error) {
	var result int
	licenseExpireDate, err := time.Parse(time.DateOnly, productLicenseResponse.ExpireDate)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not parse expireDate information. error %w", err)
	}
	duration := time.Until(licenseExpireDate)
	result = int(duration.Hours() / hoursPerDay)

	return result, err
}

func getPrettyJSON(license string) (bytes.Buffer, error) {
	var result bytes.Buffer
	err := json.Indent(&result, []byte(license), "", "\t")
	if err != nil {
		return result, wraperror.Errorf(err, "Could not parse license information.  Error %w", err)
	}

	return result, err
}

func getProductLicenseResponse(license string) (*ProductLicenseResponse, error) {
	result := &ProductLicenseResponse{}
	err := json.Unmarshal([]byte(license), result)
	if err != nil {
		return result, wraperror.Errorf(err, "Could not parse license information into structure.  Error %w", err)
	}

	return result, err
}

func returnValues(
	reportChecks []string,
	reportInfo []string,
	reportErrors []string,
	err error,
) ([]string, []string, []string, error) {
	if err != nil {
		reportErrors = append(reportErrors, err.Error())
	}

	return reportChecks, reportInfo, reportErrors, nil
}
