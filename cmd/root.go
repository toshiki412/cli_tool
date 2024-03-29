/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
)

var (
	configFile string
	setting    cfg.SettingType
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli_tool",
	Short: "",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./cli_tool.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		dir, err := file.FindCurrentDir()
		if err != nil {
			return
		}
		viper.AddConfigPath(dir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli_tool")
	}

	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	err = viper.Unmarshal(&setting)
	cobra.CheckErr(err)
}
