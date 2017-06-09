package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	appsecret string
	isColor   bool
	isVerbose bool
	isDebug   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "rai-crypto",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.PersistentFlags().StringVarP(&appsecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", color.NoColor, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")
}

func initConfig() {
}

func initColor() {
	color.NoColor = !isColor
}
