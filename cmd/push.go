/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
	"github.com/toshiki412/cli_tool/storage"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push [flags] [version id ... ]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(0), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// 引数でバージョンがない場合はlocalにあってリモートにないものをすべてアップロード
		// 引数でバージョンがある場合はそのバージョンをアップロード

		dir, err := file.DataDir()
		cobra.CheckErr(err)

		targets := make([]string, len(args))

		for i, versionId := range args {
			targets[i] = versionId
		}

		// アップロード
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				for _, versionId := range targets {
					storage.Upload(filepath.Join(dir, versionId+".zip"), versionId+".zip", conf)
				}

				// .cli_tool_localのデータをリモートの.cli_toolに同期する
			},
		})
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
