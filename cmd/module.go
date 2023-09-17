package cmd

import (
	"fmt"

	"slices"

	"github.com/spf13/viper"
)

const MODULE_TYPE = "module_type"

func readConfig(path string) error {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)     // path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %w", err)
	}
	if !viper.IsSet(MODULE_TYPE) {
		return fmt.Errorf("invalid config: %s missing", MODULE_TYPE)
	}
	module_type := viper.GetString(MODULE_TYPE)
	if !slices.Contains(MODULES_ALL, module_type) {
		return fmt.Errorf("invalid config: %s can be one of: %s", MODULE_TYPE, MODULES_ALL)
	}
	return nil
}
