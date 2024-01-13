package dump_file

import (
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
)

func Dump(dumpDir string, conf cfg.TargetFileType) {
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)

	src := filepath.Join(cwd, conf.Path)
	dest := filepath.Join(dumpDir, conf.Path)
	err = os.MkdirAll(dest, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, dest)
	cobra.CheckErr(err)
}

func Expand(dumpDir string, conf cfg.TargetFileType) {
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)

	src := filepath.Join(dumpDir, conf.Path)
	dest := filepath.Join(cwd, conf.Path)
	err = os.MkdirAll(dest, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, cwd)
	cobra.CheckErr(err)
}
