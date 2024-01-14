package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
)

const VERSION_FILE = ".cli_tool_version"
const HISTORY_FILE = ".cli_tool"
const DATADIR = ".cli_tool"

func configFiles() []string {
	return []string{"cli_tool.yaml", "cli_tool.yml"} // yamlでもymlでもいい
}

func MakeTempDir() (string, error) {
	return os.MkdirTemp("", ".cli_tool")
}

// cli_tool.yamlがあるディレクトリを探す
func FindCurrentDir() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for _, f := range configFiles() {
		return searchFile(p, f)
	}
	return "", fmt.Errorf("config file not found")
}

func searchFile(dir string, filename string) (string, error) {
	// ルートディレクトリに到達したらエラー
	if dir == filepath.Dir(dir) {
		return "", fmt.Errorf("file not found")
	}

	p := filepath.Join(dir, filename)
	_, err := os.Stat(p)
	if err != nil {
		return searchFile(filepath.Dir(dir), filename) // さらに上の階層を探す
	}
	return dir, nil
}

func ReadVersionFile() string {
	dir, err := FindCurrentDir()
	if err != nil {
		return ""
	}
	file := filepath.Join(dir, VERSION_FILE)
	data, err := readFile(file)
	if err != nil {
		return ""
	}
	return strings.Replace(data, "\n", "", -1)
}

func UpdateVersionFile(versionId string) error {
	dir, err := FindCurrentDir()
	file := filepath.Join(dir, VERSION_FILE)
	if err != nil {
		return err
	}
	return writeFile(file, versionId)
}

// データが置かれているディレクトリ(.cli_tool)を取得する
func DataDir() (string, error) {
	dir, err := FindCurrentDir()
	if err != nil {
		return "", err
	}
	d := filepath.Join(dir, DATADIR)
	s, err := os.Stat(d)
	if err != nil {
		os.Mkdir(d, os.ModePerm)
	} else if !s.IsDir() {
		return "", fmt.Errorf("datadir is file")
	}

	return d, nil
}

func UpdateHistoryFile(dir string, suffix string, newVersion cfg.VersionType) error {
	b, err := json.Marshal(newVersion)
	cobra.CheckErr(err)
	newLine := fmt.Sprintf("%s\n", string(b))

	file := filepath.Join(dir, HISTORY_FILE+suffix)
	_, err = os.Stat(file)
	if err != nil {
		writeFile(file, newLine)
		return nil
	}
	return appendFile(file, newLine)
}

func ListHistory(suffix string) []cfg.VersionType {
	dir, err := DataDir()
	cobra.CheckErr(err)

	file := filepath.Join(dir, HISTORY_FILE+suffix)
	content, err := readFile(file)
	cobra.CheckErr(err)

	lines := strings.Split(content, "\n")
	var list = make([]cfg.VersionType, 0)
	var version cfg.VersionType
	for _, line := range lines {
		if line == "" {
			continue
		}
		err = json.Unmarshal([]byte(line), &version)
		cobra.CheckErr(err)

		list = append(list, version)
	}

	return list
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(file string, data string) error {
	return os.WriteFile(file, []byte(data), os.ModePerm)
}

func appendFile(file string, data string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func findVersion(versionId string, suffix string) (cfg.VersionType, error) {
	list := ListHistory(suffix)
	for _, version := range list {
		if strings.HasPrefix(version.Id, versionId) {
			return version, nil
		}
	}
	return cfg.VersionType{}, fmt.Errorf("version not found")
}

func FindVersion(versionId string) (cfg.VersionType, error) {
	remoteVersion, err := findVersion(versionId, "")
	if err == nil {
		return remoteVersion, nil
	}

	localVersion, err := findVersion(versionId, "_local")
	if err == nil {
		return localVersion, nil
	}

	return cfg.VersionType{}, fmt.Errorf("version not found")
}

func MoveVersion(target cfg.VersionType) {
	localList := ListHistory("_local")
	remoteList := ListHistory("")

	newLocalList := make([]cfg.VersionType, 0)
	for _, version := range localList {
		if version.Id == target.Id {
			continue
		}
		newLocalList = append(newLocalList, version)
	}

	remoteList = append(remoteList, target)
	sort.Slice(remoteList, func(i, j int) bool {
		return remoteList[i].Time < remoteList[j].Time
	})

	writeFile(filepath.Join(DATADIR, HISTORY_FILE), versionListToString(remoteList))
	writeFile(filepath.Join(DATADIR, HISTORY_FILE+"_local"), versionListToString(newLocalList))
}

func versionListToString(list []cfg.VersionType) string {
	var str = ""
	for _, version := range list {
		line, err := versionToString(version)
		cobra.CheckErr(err)
		str += line
	}
	return str
}

func versionToString(version cfg.VersionType) (string, error) {
	b, err := json.Marshal(version)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n", string(b)), nil
}
