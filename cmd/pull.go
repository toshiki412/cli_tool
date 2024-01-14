/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
	"github.com/toshiki412/cli_tool/storage"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [flags] [version id]",
	Short: "pull remote version",
	Long:  `pull remote version`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// pullすると、引数が無ければ.cli_tool_versionのデータがリモートからpullされる
		// 引数があれば、そのバージョンのデータがリモートからpullされる

		// .cli_toolがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("cli_tool.yaml not found!")
			return
		}

		// 引数にversionIdがあるかどうか
		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			versionId = file.ReadVersionFile()
		}
		if versionId == "" {
			fmt.Println("version not found!")
			return
		}
		version, err := file.FindVersion(versionId)
		if err != nil {
			fmt.Println("version not found!")
			return
		}

		// 取れたバージョンをダウンロードする
		var downloadedFile string
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				downloadedFile = storage.Download(version.Id+".zip", conf)
				fmt.Println("downloaded from google storage!")
			},
		})

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(downloadedFile, filepath.Join(dataDir, versionId+".zip"))
		cobra.CheckErr(err)

		fmt.Printf("pulled successfully! version_id: %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
