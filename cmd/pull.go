/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/cfg/compress"
	"github.com/toshiki412/cli_tool/dump"
	"github.com/toshiki412/cli_tool/storage"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [flags] [version id]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")
		fmt.Println(args)

		// pullidがあるかどうか
		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			// TODO cli_tool_versionがない場合の処理
			f, err := os.Open(".cli_tool_version")
			cobra.CheckErr(err)

			data := make([]byte, 1024)
			_, err = f.Read(data)
			cobra.CheckErr(err)

			// 最新のバージョンを取得する
			versionId = strings.Replace(string(data), "\n", "", -1)

			err = f.Close()
			cobra.CheckErr(err)
		}

		// 指定のバージョンをダウンロードする
		var tmpFile string
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				tmpFile = storage.Download(fmt.Sprintf("%s.zip", versionId), conf)
				fmt.Println("downloaded from google storage!")
			},
		})

		tmpDir, err := os.MkdirTemp("", ".cli_tool")
		cobra.CheckErr(err)

		// 展開する
		compress.Decompress(tmpDir, tmpFile)

		// 展開したものを適用する
		cfg.DispatchTarget(setting.Target, cfg.TargetFuncTable{
			Mysql: func(conf cfg.TargetMysqlType) {
				dump.Import(tmpDir, conf)
			},
		})

		// .cli_tool_versionを更新する
		f, err := os.Open(".cli_tool_version")
		cobra.CheckErr(err)
		defer f.Close()
		f.WriteString(versionId)

		fmt.Printf("pulled successfully! version_id: %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
