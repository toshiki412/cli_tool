package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/dump/dump_mysql"
	"github.com/toshiki412/cli_tool/file"
)

var database string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import mysql dump file",
	Long:  `import mysql dump file`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		// cli_tool.yamlがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Printf("cli_tool.yaml not found! \n")
			fmt.Printf("please run cli_tool init\n")
			return
		}

		dumpFile := args[0]

		for _, target := range setting.Targets {
			cfg.DispatchTarget(target, cfg.TargetFuncTable{
				Mysql: func(config cfg.TargetMysqlType) {
					if database == config.Database {
						fmt.Printf("import mysql database = %s\n", config.Database)
						dump_mysql.Import(dumpFile, config)
					}
				},
				File: func(config cfg.TargetFileType) {
					// no support
				},
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&database, "database", "d", "", "database name")
}
