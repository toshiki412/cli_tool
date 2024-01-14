package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFiles(t *testing.T) {
	cwd, _ := os.Getwd()
	dir, err := FindCurrentDir()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Dir(cwd), dir)
}

// .cli_tool_versionファイルにあるバージョンを読む
func TestReadVersionFile(t *testing.T) {
	assert.Equal(t, "e60d95533ff24fcc954b216787d65bf1", ReadVersionFile())
}
