package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   NAME,
	Short: "Modular mock server for simple integration testing",
	Long:  `Modular mock server for simple integration testing`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	runCmd := NewRunCmd()
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		// no need ot print as the Execute prints it by default
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
