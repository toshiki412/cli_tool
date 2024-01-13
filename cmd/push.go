/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/cfg/compress"
	"github.com/toshiki412/cli_tool/dump"
	"github.com/toshiki412/cli_tool/storage"
)

var pushMessage string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called", setting)
		if setting.Target.Kind == "mysql" {
			// TargetMysqlConfigTypeに変換する
			var conf cfg.TargetMysqlType
			err := mapstructure.Decode(setting.Target.Config, &conf)
			cobra.CheckErr(err)

			// dbダンプ
			dumpDir, err := os.MkdirTemp("", ".cli_tool")
			cobra.CheckErr(err)
			dump.Dump(dumpDir, conf)

			// zip圧縮
			zipfile := compress.Compress(dumpDir)

			// アップロード
			var gcsConf cfg.UploadGoogleStorageType
			err = mapstructure.Decode(setting.Upload.Config, &gcsConf)
			cobra.CheckErr(err)

			_uuid, err := uuid.NewRandom()
			cobra.CheckErr(err)
			uuid := _uuid.String()
			uuid = strings.Replace(uuid, "-", "", -1)
			uuidWithZip := fmt.Sprintf("%s.zip", uuid)

			storage.Upload(zipfile, uuidWithZip, gcsConf)

			fmt.Println("pushed to google storage!")

			nowTime := time.Now()

			v := cfg.VersionType{
				Id:      uuid,
				Time:    nowTime.Unix(),
				Message: pushMessage,
			}

			b, err := json.Marshal(v)
			cobra.CheckErr(err)
			version := string(b)

			if storage.IsExist(".cli_tool", gcsConf) {
				filePath := storage.Download(".cli_tool", gcsConf)
				f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
				cobra.CheckErr(err)
				defer f.Close()

				_, err = f.WriteString(fmt.Sprintf("%s\n", version))
				cobra.CheckErr(err)
				err = f.Close()
				cobra.CheckErr(err)

				storage.Upload(filePath, ".cli_tool", gcsConf)
			} else {
				tmpDir, err := os.MkdirTemp("", ".cli_tool")
				cobra.CheckErr(err)

				tmpFile := filepath.Join(tmpDir, ".cli_tool")
				f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY, 0644)
				cobra.CheckErr(err)
				_, err = f.WriteString(fmt.Sprintf("%s\n", version))
				cobra.CheckErr(err)
				err = f.Close()
				cobra.CheckErr(err)
				storage.Upload(tmpFile, ".cli_tool", gcsConf)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// flagの追加 -m, --messageでプッシュメッセージを指定できるようにする
	pushCmd.Flags().StringVarP(&pushMessage, "message", "m", "", "message for push")
	pushCmd.MarkPersistentFlagRequired("message")
}
