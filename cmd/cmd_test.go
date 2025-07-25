package cmd_test

import (
	"os"
	"testing"

	"github.com/senzing-garage/check-self/cmd"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name"}

	cmd.Execute()
}

// func Test_Execute_completion(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "completion"}

// 	cmd.Execute()
// }

// func Test_Execute_docs(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "docs"}

// 	cmd.Execute()
// }

// func Test_Execute_help(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "--help"}

// 	cmd.Execute()
// }

// func Test_PreRun(test *testing.T) {
// 	test.Parallel()

// 	args := []string{"command-name"}
// 	cmd.PreRun(cmd.RootCmd, args)
// }

// func Test_RunE(test *testing.T) {
// 	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")

// 	err := cmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_RootCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.RootCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.RootCmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_CompletionCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.CompletionCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.CompletionCmd.RunE(cmd.CompletionCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_DocsCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.DocsCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.DocsCmd.RunE(cmd.DocsCmd, []string{})
// 	require.NoError(test, err)
// }
