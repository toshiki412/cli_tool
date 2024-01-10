/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/toshiki412/cli_tool/cfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	config     cfg.ConfigType
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli_tool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		cwd, err := os.Getwd() // カレントディレクトリの取得
		cobra.CheckErr(err)

		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli_tool")
	}

	// viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)
}
