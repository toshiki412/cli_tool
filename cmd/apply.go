/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/cfg/compress"
	"github.com/toshiki412/cli_tool/dump/dump_file"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [flags] [version id]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		var versionId = ""
		if len(args) == 1 {
			versionId = args[0]
		} else {
			// TODO cli_tool_versionがない場合の処理
			var err error = nil
			versionId, err = file.ReadVersionFile()
			cobra.CheckErr(err)
		}

		dataDir, err := file.DataDir()
		cobra.CheckErr(err)
		tmpFile := filepath.Join(dataDir, versionId+".zip")

		// .cli_tool_localを見る

		tmpDir, err := file.MakeTempDir()
		cobra.CheckErr(err)
		defer os.RemoveAll(tmpDir)

		// 展開する
		compress.Decompress(tmpDir, tmpFile)

		// 展開したものを適用する
		for _, target := range setting.Targets {
			cfg.DispatchTarget(target, cfg.TargetFuncTable{
				Mysql: func(conf cfg.TargetMysqlType) {
					dump_mysql.Import(tmpDir, conf)
				},
				File: func(conf cfg.TargetFileType) {
					dump_file.Expand(tmpDir, conf)
				},
			})
		}

		// .cli_tool_versionを更新する
		err = file.UpdateVersionFile(versionId)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
