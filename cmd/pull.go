/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

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

		// 引数にversionIdがあるかどうか
		version, err := file.GetCurrentVersion(args)
		if err != nil {
			fmt.Println("version not found!")
			return
		}

		// 取れたバージョンをダウンロードする
		var downloadedFile string
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				downloadedFile = storage.Download(version.CreateZipFile(), conf)
				fmt.Println("downloaded from google storage!")
			},
		})

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(downloadedFile, version.CreateZipFileWithDir(dataDir))
		cobra.CheckErr(err)

		fmt.Printf("pulled successfully! version_id: %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
