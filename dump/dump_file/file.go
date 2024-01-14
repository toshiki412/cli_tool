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
	s, err := os.Stat(src)
	cobra.CheckErr(err)

	dest := filepath.Join(dumpDir, conf.Path)
	destDir := dest
	if !s.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, dest)
	cobra.CheckErr(err)
}

func Expand(dumpDir string, conf cfg.TargetFileType) {
	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)

	src := filepath.Join(dumpDir, conf.Path)
	s, err := os.Stat(src)
	cobra.CheckErr(err)

	dest := filepath.Join(cwd, conf.Path)
	destDir := dest
	if !s.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, dest)
	cobra.CheckErr(err)
}
