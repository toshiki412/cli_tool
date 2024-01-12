/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
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
			dumpfile := processMysqlDump(conf)
			fmt.Println(dumpfile)

			// アップロード
			uploadGoogleStorage()
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}

func processMysqlDump(conf cfg.TargetMysqlConfigType) string {
	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.DBName = conf.Database

	dumpDir, err := os.MkdirTemp("", ".cli_tool")
	cobra.CheckErr(err)

	dumpFileNameFormat := fmt.Sprintf("%s-20060102150405", conf.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFileNameFormat, config.DBName)
	cobra.CheckErr(err)

	// dump database to file
	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Printf("successfully dumped to file %s\n", dumpFileNameFormat)

	dumper.Close()

	return filepath.Join(dumpDir, dumpFileNameFormat+".sql")
}

func uploadGoogleStorage() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	f, err := os.Open("Readme.md")
	cobra.CheckErr(err)
	defer f.Close()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	// defer cancel()

	o := client.Bucket("clitoolbacket0001").Object("Readme.md")
	// o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)
	err = wc.Close()
	cobra.CheckErr(err)

}
