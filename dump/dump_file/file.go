package dump_file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
	"github.com/toshiki412/cli_tool/file"
)

func Dump(dumpDir string, conf cfg.TargetFileType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" file dump ... (path: %s)\n", conf.Path)
	s.Start()

	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)

	src := filepath.Join(cwd, conf.Path)
	stat, err := os.Stat(src)
	cobra.CheckErr(err)

	dest := filepath.Join(dumpDir, conf.Path)
	destDir := dest
	if !stat.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, dest)
	cobra.CheckErr(err)

	s.Stop()
	fmt.Printf("✔︎ file dump completed. (path: %s)\n", conf.Path)
}

func Expand(dumpDir string, conf cfg.TargetFileType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" restore file(s) ... (path: %s)\n", conf.Path)
	s.Start()

	cwd, err := file.FindCurrentDir()
	cobra.CheckErr(err)

	src := filepath.Join(dumpDir, conf.Path)
	stat, err := os.Stat(src)
	cobra.CheckErr(err)

	dest := filepath.Join(cwd, conf.Path)
	destDir := dest
	if !stat.IsDir() {
		destDir = filepath.Dir(dest)
	}
	err = os.MkdirAll(destDir, os.ModePerm)
	cobra.CheckErr(err)

	err = cp.Copy(src, dest)
	cobra.CheckErr(err)

	s.Stop()
	fmt.Printf("✔︎ restore file(s) completed. (path: %s)\n", conf.Path)
}
