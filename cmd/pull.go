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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		// 引数にversionIdがあるかどうか
		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			var err error = nil
			versionId, err = file.ReadVersionFile()
			cobra.CheckErr(err)
		}

		// 指定のバージョンをダウンロードする
		var tmpFile string
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				tmpFile = storage.Download(versionId+".zip", conf)
				fmt.Println("downloaded from google storage!")
			},
		})

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		err = os.Rename(tmpFile, filepath.Join(dir, versionId+".zip"))
		cobra.CheckErr(err)

		fmt.Printf("pulled successfully! version_id: %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
