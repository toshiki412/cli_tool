/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		// .cli_tool.yamlがあるかどうか
		dir, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("cli_tool.yaml not found!")
			return
		}

		viper.AddConfigPath(dir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli_tool")
		fmt.Println(dir)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&setting)
	cobra.CheckErr(err)
}
