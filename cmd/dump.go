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
	"github.com/toshiki412/cli_tool/compress"
	"github.com/toshiki412/cli_tool/dump/dump_file"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
)

var dumpMessage string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump current data",
	Long:  `dump current data`,
	Run: func(cmd *cobra.Command, args []string) {
		// dumpすると.cli_tool下にzipファイルが生成され、
		// .cli_tool_localに新しいバージョンの履歴が追加され、
		// .cli_tool_versionが更新される

		// .cli_toolがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("cli_tool.yaml not found!")
			return
		}

		// dbダンプ
		// dumpDirにダンプしたデータが入る
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
		versionId := _uuid.String() // 新しいバージョンのuuidを振る
		versionId = strings.Replace(versionId, "-", "", -1)

		// .cli_toolに持っていく
		dir, err := file.DataDir()
		cobra.CheckErr(err)
		dest := filepath.Join(dir, versionId+".zip") // .cli_tool/xxxxxxxx.zip
		err = os.Rename(zipfile, dest)
		cobra.CheckErr(err)

		// 新しいバージョンのstructを作成
		newVersion := cfg.VersionType{
			Id:      versionId,
			Time:    time.Now().Unix(),
			Message: dumpMessage,
		}

		// 生成した新しいバージョンをローカルに保存する
		file.UpdateHistoryFile(dir, "_local", newVersion)

		// .cli_tool_versionをdumpしたバージョンに更新する
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
