package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove dump",
	Long:  `remove dump`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {

		// cli_tool.yamlがあるかどうか
		_, err := file.FindCurrentDir()
		if err != nil {
			fmt.Printf("cli_tool.yaml not found! \n")
			fmt.Printf("please run cli_tool init\n")
			return
		}

		versionId := args[0]

		ds := file.ReadLocalDataFile()

		newHistories := []cfg.VersionType{}
		for _, v := range ds.Histories {
			if v.Id != versionId {
				newHistories = append(newHistories, v)
			}
		}
		if len(ds.Histories) == len(newHistories) {
			fmt.Printf("version not found. %s\n", versionId)
			return
		}
		ds.Histories = newHistories
		file.WriteLocalDataFile(ds)

		dir, err := file.DataDir()
		cobra.CheckErr(err)
		os.Remove(filepath.Join(dir, versionId))

		fmt.Printf("removed. %s\n", versionId)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
