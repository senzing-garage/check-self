/*
 */
package cmd

import (
	"context"
	"os"

	"github.com/senzing-garage/check-self/checkself"
	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Check the environment in which senzing-tool runs"
	Use   string = "check-self"
	Long  string = `
check-self long description.
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.ConfigPath,
	option.Configuration,
	option.DatabaseURL,
	option.EngineSettings,
	option.EngineLogLevel,
	option.GrpcURL,
	option.InputURL,
	option.LicenseDaysLeft,
	option.LicenseRecordsPercent,
	option.LicenseStringBase64,
	option.LogLevel,
	option.ObserverURL,
	option.ResourcePath,
	option.SenzingDirectory,
	option.SupportPath,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command.
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command.
func RunE(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	checkSelf := &checkself.BasicCheckSelf{
		ConfigPath:                 viper.GetString(option.ConfigPath.Arg),
		DatabaseURL:                viper.GetString(option.DatabaseURL.Arg),
		Settings:                   viper.GetString(option.EngineSettings.Arg),
		EngineLogLevel:             viper.GetString(option.EngineLogLevel.Arg),
		ErrorLicenseDaysLeft:       viper.GetString(option.LicenseDaysLeft.Arg),
		ErrorLicenseRecordsPercent: viper.GetString(option.LicenseRecordsPercent.Arg),
		GrpcURL:                    viper.GetString(option.GrpcPort.Arg),
		InputURL:                   viper.GetString(option.InputURL.Arg),
		LicenseStringBase64:        viper.GetString(option.LicenseStringBase64.Arg),
		LogLevel:                   viper.GetString(option.LogLevel.Arg),
		ObserverURL:                viper.GetString(option.ObserverGrpcPort.Arg),
		ResourcePath:               viper.GetString(option.ResourcePath.Arg),
		SenzingDirectory:           viper.GetString(option.SenzingDirectory.Arg),
		SupportPath:                viper.GetString(option.SupportPath.Arg),
	}
	return checkSelf.CheckSelf(ctx)
}

// Used in construction of cobra.Command.
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}
