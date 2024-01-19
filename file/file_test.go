package file

import (
	"fmt"
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

func createFile(home string, file string, content string) {
	os.WriteFile(filepath.Join(home, file), []byte(content), os.ModePerm)
}

func TestFindCurrentDir(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	dir, err := FindCurrentDir()
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("file not found"), err)
	assert.Equal(t, "", dir)

	createFile(home, "cli_tool.yaml", "")

	dir, err = FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, home, dir)
}

// .cli_tool_versionファイルにあるバージョンを読む
func TestReadVersionFile(t *testing.T) {
	home := setup(t)
	defer teardown(home)

	assert.Equal(t, "", ReadVersionFile())

	createFile(home, ".cli_tool_version", "")
	assert.Equal(t, "", ReadVersionFile())

	createFile(home, "cli_tool.yaml", "")
	assert.Equal(t, "", ReadVersionFile())

	createFile(home, ".cli_tool_version", "1.2.3")
	assert.Equal(t, "1.2.3", ReadVersionFile())

	createFile(home, ".cli_tool_version", "1.2.3\n")
	assert.Equal(t, "1.2.3", ReadVersionFile())
}

func TestUppdateVersionFile(t *testing.T) {
	home := setup(t)
	defer teardown(home)
	createFile(home, "cli_tool.yaml", "")

	err := UpdateVersionFile("1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3", ReadVersionFile())

	err = UpdateVersionFile("4.5.6\n")
	assert.NoError(t, err)
	assert.Equal(t, "4.5.6", ReadVersionFile())
}
