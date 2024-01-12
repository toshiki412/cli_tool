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
	"github.com/toshiki412/cli_tool/cfg/compress"
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
		fmt.Println("push called", setting)
		if setting.Target.Kind == "mysql" {
			// TargetMysqlConfigTypeに変換する
			var conf cfg.TargetMysqlType
			err := mapstructure.Decode(setting.Target.Config, &conf)
			cobra.CheckErr(err)

			// dbダンプ
			dumpDir, err := os.MkdirTemp("", ".cli_tool")
			cobra.CheckErr(err)
			processMysqlDump(dumpDir, conf)

			// zip圧縮
			zipfile := compress.Compress(dumpDir)

			// アップロード
			var gcsConf cfg.UploadGoogleStorageType
			err = mapstructure.Decode(setting.Upload.Config, &gcsConf)
			cobra.CheckErr(err)
			uploadGoogleStorage(zipfile, gcsConf)

			fmt.Println("pushed to google storage!")
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}

func processMysqlDump(dumpDir string, conf cfg.TargetMysqlType) {
	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.DBName = conf.Database

	dumpFileNameFormat := fmt.Sprintf("%s-%s", "mysql", conf.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFileNameFormat, config.DBName)
	cobra.CheckErr(err)

	// dump database to file
	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Println("successfully dumped to file.")

	dumper.Close()
}

func uploadGoogleStorage(target string, conf cfg.UploadGoogleStorageType) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	f, err := os.Open(target)
	cobra.CheckErr(err)
	defer f.Close()

	var uploadpath = ""
	if conf.Dir == "" {
		uploadpath = "upload.zip"
	} else {
		uploadpath = filepath.Join(conf.Dir, "upload.zip")
	}

	o := client.Bucket(conf.Bucket).Object(uploadpath)

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)
	err = wc.Close()
	cobra.CheckErr(err)
}
