/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
	"github.com/toshiki412/cli_tool/storage"
)

var remote bool

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list history",
	Long:  `list history`,
	Run: func(cmd *cobra.Command, args []string) {
		// lsはlocalのバージョン履歴を表示する
		// -rオプションがある場合はリモートのバージョン履歴も表示する

		// .cli_toolがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("cli_tool.yaml not found!")
			return
		}

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)

		// リモートのバージョン履歴を読み込む
		var remoteList []cfg.VersionType
		if remote {
			var tmpFile string
			cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
				Gcs: func(conf cfg.StorageGoogleStorageType) {
					tmpFile = storage.Download(".cli_tool", conf)
				},
			})

			os.Rename(tmpFile, filepath.Join(dataDir, ".cli_tool"))
			remoteList = file.ListHistory("")
		}
		// ローカルのバージョン履歴を読み込む
		localList := file.ListHistory("_local")

		// 表示する
		if remote {
			fmt.Println("~~~ remote history ~~~")
			fmt.Println("id\ttime\tmessage")
			for _, version := range remoteList {
				d := time.Unix(version.Time, 0).Format("2006-01-02 15:04:05")
				fmt.Printf("%s\t%s\t%s\n", version.Id, d, version.Message)
			}
		}
		fmt.Println("~~~ local history ~~~")
		fmt.Println("id\ttime\tmessage")
		for _, version := range localList {
			d := time.Unix(version.Time, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("%s\t%s\t%s\n", version.Id, d, version.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&remote, "remote", "r", false, "show with remote history")
}
