/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/compress"
	"github.com/toshiki412/cli_tool/dump/dump_file"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [flags] [version id]",
	Short: "apply version",
	Long:  `apply version`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// applyすると、.cli_tool下のzipファイルが展開され、
		// .cli_tool_localに新しいバージョンの履歴が追加され、
		// .cli_tool_versionがapplyしたバージョンに更新される

		// .cli_toolがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Println("cli_tool.yaml not found!")
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
		tmpFile := version.CreateZipFileWithDir(dataDir)

		// ダウンロードしてない場合の処理
		s, err := os.Stat(tmpFile)
		if err != nil || s.IsDir() {
			fmt.Printf("file not found.\nplease run cli_tool pull %s\n", version.Id)
		}

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

		// .cli_tool_versionをapplyしたバージョンに更新する
		err = file.UpdateVersionFile(version.Id)
		cobra.CheckErr(err)

		fmt.Printf("applied successfully! version_id: %s\n", version.Id)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
