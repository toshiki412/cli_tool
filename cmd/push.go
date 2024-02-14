/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
	"github.com/toshiki412/cli_tool/storage"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push [flags] [version id ... ]",
	Short: "upload version to remote storage",
	Long:  `upload version to remote storage`,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(0), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// pushすると、引数が無ければ.cli_tool_versionのデータがリモートにpushされる
		// 引数があれば、そのバージョンのデータがリモートにpushされる
		// pushすると.cli_tool_localから履歴がなくなりリモートに移動する

		// cli_tool.yamlがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Printf("cli_tool.yaml not found! \n")
			fmt.Printf("please run cli_tool init\n")
			return
		}

		// 引数にversionIdがあるかどうか
		version, err := file.GetCurrentVersion(args)
		if err != nil {
			fmt.Println("version not found!")
			return
		}

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)

		// アップロード
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				storage.Upload(version.CreateZipFileWithDir(dataDir), version.CreateZipFile(), conf)

				// FIXME .cli_tool_localのデータをリモートの.cli_toolに同期する
				file.MoveVersionToRemote(version)
				storage.Upload(filepath.Join(dataDir, ".cli_tool"), ".cli_tool", conf)
			},
		})
		fmt.Printf("push successfully! version_id: %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
