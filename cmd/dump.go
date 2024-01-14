/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/cfg/compress"
	"github.com/toshiki412/cli_tool/dump/dump_file"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
)

var dumpMessage string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
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

		// .cli_toolに移動
		dir, err := file.DataDir()
		cobra.CheckErr(err)
		dest := filepath.Join(dir, fmt.Sprintf("%s.zip", versionId))
		err = os.Rename(zipfile, dest)
		cobra.CheckErr(err)

		// .cli_tool/.cli_tool_local この中がローカル
		// .cli_tool/.cli_tool(-remote) これがリモート

		newVersion := cfg.VersionType{
			Id:      versionId,
			Time:    time.Now().Unix(),
			Message: dumpMessage,
		}

		file.UpdateHistoryFile(dir, "_local", newVersion)
		file.UpdateVersionFile(versionId)

		fmt.Printf("dump success! version id: %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	// flagの追加 -m, --messageでプッシュメッセージを指定できるようにする
	dumpCmd.Flags().StringVarP(&dumpMessage, "message", "m", "", "message for push")
	dumpCmd.MarkFlagRequired("message")
}
