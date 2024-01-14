/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/cfg/compress"
	"github.com/toshiki412/cli_tool/dump/dump_file"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
	"github.com/toshiki412/cli_tool/storage"
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
		// dbダンプ
		dumpDir, err := file.MakeTempDir()
		cobra.CheckErr(err)
		defer os.RemoveAll(dumpDir)

		for _, target := range setting.Targets {
			cfg.DispatchTarget(target, cfg.TargetFuncTable{
				Mysql: func(conf cfg.TargetMysqlType) {
					dump_mysql.Dump(dumpDir, conf)
				},
				File: func(conf cfg.TargetFileType) {
					dump_file.Dump(dumpDir, conf)
				},
			})
		}

		// zip圧縮
		zipfile := compress.Compress(dumpDir)

		_uuid, err := uuid.NewRandom()
		cobra.CheckErr(err)
		versionId := _uuid.String()
		versionId = strings.Replace(versionId, "-", "", -1)

		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				// アップロード
				storage.Upload(zipfile, fmt.Sprintf("%s.zip", versionId), conf)
			},
		})
		// TODO .cli_tool_versionを更新する
		fmt.Printf("pushed successfully! version_id: %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
