/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

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

		// cli_tool.yamlがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Printf("cli_tool.yaml not found! \n")
			fmt.Printf("please run cli_tool init\n")
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
					dump_mysql.Dump(dump_mysql.MysqlDumpFile(dumpDir, conf), conf)
				},
				File: func(conf cfg.TargetFileType) {
					dump_file.Dump(dumpDir, conf)
				},
			})
		}

		// zip圧縮
		zipfile := compress.Compress(dumpDir)

		// 新しいバージョンのstructを作成
		versionId, err := file.NewUUID()
		cobra.CheckErr(err)
		newVersion := cfg.VersionType{
			Id:      versionId,
			Time:    time.Now().Unix(),
			Message: dumpMessage,
		}

		// .cli_toolに持っていく
		dir, err := file.DataDir()
		cobra.CheckErr(err)
		dest := newVersion.CreateZipFileWithDir(dir) // .cli_tool/xxxxxxxx.zip
		err = file.MoveFile(zipfile, dest)
		cobra.CheckErr(err)

		local := file.ReadLocalDataFile()
		local.Histories = append(local.Histories, newVersion)

		// 生成した新しいバージョンをローカルに保存する
		file.WriteLocalDataFile(local)

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
