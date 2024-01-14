package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) string {
	dir, err := os.MkdirTemp("", ".cli_tool_test")
	assert.NoError(t, err)
	os.Chdir(dir)
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	return cwd
}

func teardown(dir string) {
	os.RemoveAll(dir)
}

func TestFindCurrentDir(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	dir, err := FindCurrentDir()
	assert.Error(t, err)
	assert.Equal(t, "", dir)

	os.WriteFile(filepath.Join(home, "cli_tool.yaml"), []byte(""), os.ModePerm)

	dir, err = FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, home, dir)
}

// .cli_tool_versionファイルにあるバージョンを読む
// func TestReadVersionFile(t *testing.T) {
// 	assert.Equal(t, "b1084e3432394fec9425e71e21a43616", ReadVersionFile())
// }
