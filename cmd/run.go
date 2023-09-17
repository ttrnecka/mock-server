/*
Copyright © 2023 Tomas Trnecka tomas.trnecka@veritas.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
func NewRunCmd() *cobra.Command {

	var runCmd = &cobra.Command{
		Use:          "run",
		Short:        "Run mock server",
		Long:         `Run mock server on specified port`,
		RunE:         run,
		SilenceUsage: true,
	}
	runCmd.Flags().IntP("port", "p", 2222, "Port to run the mock server on")
	viper.BindPFlag("port", runCmd.Flags().Lookup("port"))
	runCmd.Flags().StringP("module", "m", "", "Module to mock. Module needs to exists in modules directory")
	runCmd.MarkFlagRequired("module")
	viper.BindPFlag("module", runCmd.Flags().Lookup("module"))

	return runCmd
}

func run(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(MODULES_FOLDER); os.IsNotExist(err) {
		return fmt.Errorf("Directory %s does not exist.", MODULES_FOLDER)
	}
	module_folder := fmt.Sprintf("%s/%s", MODULES_FOLDER, viper.GetString("module"))
	if _, err := os.Stat(module_folder); os.IsNotExist(err) {
		return fmt.Errorf("Directory %s does not exist.", module_folder)
	}
	err := readConfig(module_folder)
	if err != nil {
		return err
	}
	module_type := viper.GetString("module_type")

	if module_type == MODULE_SSH {
		return sshServer(viper.GetInt("port"))
	}
	return fmt.Errorf("module %s not implemented yet", module_type)
}
