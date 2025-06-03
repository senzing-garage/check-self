package checkself_test

import (
	"context"
	"testing"

	"github.com/senzing-garage/check-self/checkself"
	"github.com/senzing-garage/go-helpers/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	badJSON               = "}{"
	db2URL                = "db2://username:password@database/?schema=schemaname"
	expectedQuestionMarks = "????"
	mssqlURL              = "mssql://username:password@server:port:database/?driver=mssqldriver"
	mysqlURL              = "mysql://username:password@hostname:3306/?schema=schemaname"
	ociURL                = "oci://username:password@database"
	postgresqlURL         = "postgresql://username:password@hostname:5432:database/?schema=schemaname:"
	sqlite3URL            = "sqlite3://na:na@/tmp/sqlite/G2C.db"
	variableName          = "VariableName"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

// func TestBasicCheckSelf_CheckSelf_Null(test *testing.T) {
// 	ctx := test.Context()
// 	testObject := &BasicCheckSelf{}
// 	err := testObject.CheckSelf(ctx)
// 	require.NoError(test, err)
// }

func TestBasicCheckSelf_Break(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	newReportChecks, newReportInfo, newReportErrors, err := testObject.Break(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Empty(test, newReportChecks)
	assert.Empty(test, newReportInfo)
	assert.Empty(test, newReportErrors)
}

func TestBasicCheckSelf_Break_badReportErrors(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	badReportErrors := []string{"example error text"}
	reportChecks, reportInfo, reportErrors, err := testObject.Break(ctx, reportChecks(), reportInfo(), badReportErrors)
	require.Error(test, err)
	require.Empty(test, reportChecks, "reportChecks")
	require.Empty(test, reportInfo, "reportInfo")
	require.NotEmpty(test, reportErrors, "reportErrors")
}

func TestBasicCheckSelf_CheckDatabaseSchema(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckDatabaseSchema(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Empty(test, newReportChecks)
	assert.Empty(test, newReportInfo)
	assert.Empty(test, newReportErrors)
}

func TestBasicCheckSelf_CheckDatabaseSchema_badDatabaseURL(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := `SENZING_TOOLS_DATABASE_URL = bad-database-URL is misconfigured. Could not create a database connector. For more information, visit https://hub.senzing.com/...  Error: {"function": "connector.NewConnector", "error": {"function": "connector.NewConnector", "text": "unknown database scheme: ", "error": "connector"}}`
	testObject := getTestObject(ctx, test)
	badReportErrors := []string{}
	testObject.DatabaseURL = "bad-database-URL"
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

func TestBasicCheckSelf_CheckLicense(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckLicense(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	assert.Len(test, newReportInfo, 1)
	assert.Empty(test, newReportErrors)
}

func TestBasicCheckSelf_CheckLicense_badGetDatabaseURL(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
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
                "BACKEND": "SQL",
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
}

func TestBasicCheckSelf_CheckSelf(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	err := testObject.CheckSelf(ctx)
	require.NoError(test, err)
}

func TestBasicCheckSelf_CheckSelf_badSettings(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	testObject.Settings = `{}`
	err := testObject.CheckSelf(ctx)
	require.Error(test, err)
}

func TestBasicCheckSelf_CheckSenzingConfiguration_badGetDefaultConfigID(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	// expected := expectedQuestionMarks
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
                "BACKEND": "SQL",
                "CONNECTION": "sqlite3://na:na@/tmp/sqlite/G2C-empty.db"
            }
        }
        `
	testObject.DatabaseURL = "sqlite3://na:na@/tmp/sqlite/G2C-empty.db"
	newReportChecks, newReportInfo, _, err := testObject.CheckSenzingConfiguration(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	// assert.Len(test, newReportErrors, 1)
	assert.Empty(test, newReportInfo)
}

func TestBasicCheckSelf_CheckSettings(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	testObject := getTestObject(ctx, test)
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckSettings(
		ctx,
		reportChecks(),
		reportInfo(),
		reportErrors(),
	)
	require.NoError(test, err)
	assert.Len(test, newReportChecks, 1)
	assert.Empty(test, newReportInfo)
	assert.Empty(test, newReportErrors)
}

func TestBasicCheckSelf_CheckSettings_badSettings(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := `SENZING_TOOLS_ENGINE_SETTINGS - {"function": "settingsparser.New", "text": "incorrect JSON syntax in }{", "error": "settingsparser"}`
	testObject := getTestObject(ctx, test)
	testObject.Settings = badJSON
	newReportChecks, newReportInfo, newReportErrors, err := testObject.CheckSettings(
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

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

// func TestBasicCheckSelf_CheckSettings_buildAndCheckSettingsBreak_badReportErrors(test *testing.T) {
// 	ctx := test.Context()
// 	testObject := getTestObject(ctx, test)
// 	badReportErrors := []string{"example error text"}
// 	_, _, _, err := testObject.Break(ctx, reportChecks(), reportInfo(), badReportErrors)
// 	require.Error(test, err)
// }

// func TestBasicCheckSelf_CheckSettings_getDatabaseURL_badSettingsLength(test *testing.T) {
// 	ctx := test.Context()
// 	testObject := getTestObject(ctx, test)
// 	testObject.Settings = ""
// 	newReportChecks, newReportInfo, newReportErrors, err := testObject.Break(ctx, reportChecks(), reportInfo(), reportErrors())
// 	require.NoError(test, err)
// 	assert.Len(test, newReportChecks, 0)
// 	assert.Len(test, newReportInfo, 0)
// 	assert.Len(test, newReportErrors, 0)
// }

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func TestBasicCheckSelf_CheckDatabaseURL_sqlite3(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, sqlite3URL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_postgresql(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, postgresqlURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_mysql(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, mysqlURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_db2(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, db2URL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_oci(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, ociURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_mssql(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	actual := checkself.CheckDatabaseURL(ctx, variableName, mssqlURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_badURLParse(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "VariableName = \n\tnot-a-URL is misconfigured. Could not parse database URL. For more information, visit https://hub.senzing.com/...  Error: parse \"\\n\\tnot-a-URL\": net/url: invalid control character in URL"
	badDatabaseURL := "\n\tnot-a-URL"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	assert.Equal(test, expected, actual[0])
}

func TestBasicCheckSelf_CheckDatabaseURL_badURLParse_postgres(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	badDatabaseURL := "postgresql://username:password@hostname:5432:database/?schema=schemaname"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_badSqliteURL(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := expectedQuestionMarks
	badDatabaseURL := "sqlite3://na:na@host.com:port//tmp/nodatabase.db"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	sink(expected, actual)
}

func TestBasicCheckSelf_CheckDatabaseURL_badSchemaLength(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "VariableName = not-a-URL is misconfigured. A database scheme is needed (e.g. postgresql://...). For more information, visit https://hub.senzing.com/..."
	badDatabaseURL := "not-a-URL"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	assert.Equal(test, expected, actual[0])
}

func TestBasicCheckSelf_CheckDatabaseURL_badSchema(test *testing.T) {
	test.Parallel()
	ctx := test.Context()
	expected := "VariableName = badScheme://xxx is misconfigured. Scheme 'badscheme://' is not recognized. For more information, visit https://hub.senzing.com/..."
	badDatabaseURL := "badScheme://xxx"
	actual := checkself.CheckDatabaseURL(ctx, variableName, badDatabaseURL)
	assert.Equal(test, expected, actual[0])
}

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *checkself.BasicCheckSelf {
	t.Helper()

	_ = ctx
	settings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &checkself.BasicCheckSelf{
		Settings: settings,
	}

	return result
}

func reportChecks() []string {
	return []string{}
}

func reportErrors() []string {
	return []string{}
}

func reportInfo() []string {
	return []string{}
}

func sink(x string, y []string) {
	_ = x
	_ = y
}
