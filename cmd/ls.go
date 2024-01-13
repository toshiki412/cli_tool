/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/storage"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// .cli_toolをstorageから取得する
		var tmpFile string
		cfg.DispatchStorages(setting.Storage, cfg.StorageFuncTable{
			Gcs: func(conf cfg.StorageGoogleStorageType) {
				tmpFile = storage.Download(".cli_tool", conf)
			},
		})

		// 読み込む
		data, err := os.ReadFile(tmpFile)
		cobra.CheckErr(err)
		fmt.Println(string(data))
		lines := strings.Split(string(data), "\n")

		// 表示する
		fmt.Println("id\ttime\tmessage")
		for _, line := range lines {
			if line == "" {
				continue
			}
			var version cfg.VersionType
			err := json.Unmarshal([]byte(line), &version)
			cobra.CheckErr(err)
			d := time.Unix(version.Time, 0).Format("2006-01-02 15:04:05")

			fmt.Printf("%s\t%s\t%s\n", version.Id, d, version.Message)
		}

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
