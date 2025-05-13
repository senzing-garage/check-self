//go:build linux

package checkself_test

import (
	"testing"

	"github.com/senzing-garage/check-self/checkself"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicCheckSelf_CheckSelf_Paths(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	senzingPath := settings.GetSenzingPath()
	testObject := &checkself.BasicCheckSelf{
		ConfigPath:   "/etc/opt/senzing",
		DatabaseURL:  "sqlite3://na:na@/tmp/sqlite/G2C.db",
		ResourcePath: senzingPath + "/er/resources",
		SupportPath:  senzingPath + "/data",
	}
	err := testObject.CheckSelf(ctx)
	require.NoError(test, err)
}

func TestBasicCheckSelf_CheckDatabaseSchema_noSchemaInstalled(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "Senzing database schema has not been installed in sqlite3://na:na@/tmp/sqlite/G2C-empty.db. For more information, visit https://hub.senzing.com/...  Error: checker.IsSchemaInstalled.row.Scan error: no such table: DSRC_RECORD"
	testObject := getTestObject(ctx, test)
	badReportErrors := []string{}
	testObject.DatabaseURL = "sqlite3://na:na@/tmp/sqlite/G2C-empty.db"
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckDatabaseSchema(
		ctx,
		reportChecks(),
		reportInfo(),
		badReportErrors,
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	assert.Empty(test, newReportInfo)
	assert.Len(test, newReportErrors, 1)
	assert.Equal(test, expected, newReportErrors[0])
}

func TestBasicCheckSelf_CheckLicense_badGetLicense(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "Could not get count of records.  Error checker.RecordCount.row.Scan error: no such table: DSRC_RECORD"
	testObject := getTestObject(ctx, test)
	testObject.Settings = `
        {
            "PIPELINE": {
                "CONFIGPATH": "/etc/opt/senzing",
                "LICENSESTRINGBASE64": "badLicense",
                "RESOURCEPATH": "/opt/senzing/er/resources",
                "SUPPORTPATH": "/opt/senzing/data"
            },
            "SQL": {
                "BACKEND": "SQL",
                "CONNECTION": "sqlite3://na:na@/tmp/sqlite/G2C-empty.db"
            }
        }
        `
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckLicense(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	assert.Empty(test, newReportInfo)
	assert.Len(test, newReportErrors, 1)
	assert.Equal(test, expected, newReportErrors[0])
}

func TestBasicCheckSelf_CheckSettings_badDatabaseURLs(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	// expected := "????"
	testObject := getTestObject(ctx, test)
	testObject.Settings = `
        {
            "PIPELINE": {
                "CONFIGPATH": "/etc/opt/senzing",
                "LICENSESTRINGBASE64": "",
                "RESOURCEPATH": "/opt/senzing/er/resources",
                "SUPPORTPATH": "/opt/senzing/data"
            },
            "SQL": {
                "BACKEND": "SQL"
            }
        }
        `
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckSettings(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	assert.Empty(test, newReportInfo)
	// assert.Equal(test, expected, newReportErrors[0])
	assert.Len(test, newReportErrors, 2)
}

func TestBasicCheckSelf_CheckDatabaseURL_badSqliteURL_stat(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "VariableName = sqlite3://na:na@/tmp/nodatabase.db is misconfigured. Could not find /tmp/nodatabase.db. For more information, visit https://hub.senzing.com/..."
	badDatabaseURL := "sqlite3://na:na@/tmp/nodatabase.db"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	assert.Equal(test, expected, actual[0])
}
