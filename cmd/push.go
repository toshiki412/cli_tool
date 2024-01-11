/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JamesStewy/go-mysqldump"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called", config)
		if config.Target.Kind == "mysql" {
			// TargetMysqlConfigTypeに変換する
			var conf cfg.TargetMysqlConfigType
			err := mapstructure.Decode(config.Target.Config, &conf)
			cobra.CheckErr(err)
			fmt.Println(conf)
			// dbダンプ
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}

func processMysqlDump(conf cfg.TargetMysqlConfigType) string {
	dumpDir, err := os.MkdirTemp("", ".cli_tool")
	dumpFileNameFormat := fmt.Sprintf("%s-20060102150405", conf.Database)

	dns := fmt.Sprintf("%s:%s@%s:%d/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := sql.Open("mysql", dns)
	cobra.CheckErr(err)

	// register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFileNameFormat)
	cobra.CheckErr(err)

	// dump database to file
	resultFileName, err := dumper.Dump()
	cobra.CheckErr(err)

	fmt.Printf("successfully dumped to file %s\n", resultFileName)

	dumper.Close()

	return resultFileName
}
