/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullId string

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
		// var versionId = ""
		// if len(args) == 1 {
		// 	versionId = args[0]
		// } else {
		// 	// 最新のバージョンを取得する
		// 	versionId = "latest" // TODO
		// }

		// 指定のバージョンをダウンロードする
		// if setting.Upload.Kind == "gcs" {
		// 	var gcsConf cfg.UploadGoogleStorageType
		// 	err := mapstructure.Decode(setting.Upload.Config, &gcsConf)
		// 	cobra.CheckErr(err)

		// 	tmpFile := storage.Download(fmt.Sprintf("%s.zip", versionId), gcsConf)
		// }

		// 展開する

		// 展開したものを適用する
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVar(&pullId, "id", "", "version id") // 指定したIDのバージョンを取得する

}
