package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
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
	if err != nil {
		return err
	}
	file := filepath.Join(dir, VERSION_FILE)
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

func findVersion(versionId string, suffix string) (cfg.VersionType, error) {
	ds := readDataFile(suffix)
	for _, version := range ds.Histories {
		if strings.HasPrefix(version.Id, versionId) {
			return version, nil
		}
	}
	return cfg.VersionType{}, fmt.Errorf("version not found")
}

func findVersionLocalAndRemote(versionId string) (cfg.VersionType, error) {
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

func ReadLocalDataFile() cfg.DataType {
	return readDataFile("_local")
}

func ReadRemoteDataFile() cfg.DataType {
	return readDataFile("")
}

func readDataFile(suffix string) (ds cfg.DataType) {
	ds = cfg.DataType{
		Version:   "1",
		Histories: []cfg.VersionType{},
	}

	dir, err := DataDir()
	if err != nil {
		return
	}
	file := filepath.Join(dir, HISTORY_FILE+suffix)
	content, err := readFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(content), &ds)
	cobra.CheckErr(err)
	return ds
}

func MoveVersionToRemote(version cfg.VersionType) {
	local := ReadLocalDataFile()
	remote := ReadRemoteDataFile()

	newLocalList := make([]cfg.VersionType, 0)
	for _, ver := range local.Histories {
		if ver.Id == version.Id {
			continue
		}
		newLocalList = append(newLocalList, ver)
	}

	local.Histories = newLocalList

	remote.Histories = append(remote.Histories, version)
	sort.Slice(remote.Histories, func(i, j int) bool {
		return remote.Histories[i].Time < remote.Histories[j].Time
	})

	err := WriteLocalDataFile(local)
	cobra.CheckErr(err)
	err = WriteRemoteDataFile(remote)
	cobra.CheckErr(err)
}

func WriteLocalDataFile(d cfg.DataType) error {
	return writeDataFile(d, "_local")
}

func WriteRemoteDataFile(d cfg.DataType) error {
	return writeDataFile(d, "")
}

func writeDataFile(d cfg.DataType, suffix string) error {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}
	dir, err := DataDir()
	if err != nil {
		return err
	}
	err = writeFile(filepath.Join(dir, HISTORY_FILE+suffix), string(b))
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentVersion(args []string) (cfg.VersionType, error) {
	// 引数にversionIdがあるかどうか
	var versionId = ""
	if len(args) == 1 {
		versionId = args[0]
	} else {
		// 引数がない場合は.cli_tool_versionを見る
		versionId = ReadVersionFile()
	}

	if versionId == "" {
		return cfg.VersionType{}, fmt.Errorf("version not found")
	}

	version, err := findVersionLocalAndRemote(versionId)
	if err != nil {
		return cfg.VersionType{}, fmt.Errorf("version not found")
	}

	return version, nil
}

func NewUUID() (string, error) {
	_uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	uuid := _uuid.String() // 新しいバージョンのuuidを振る
	uuid = strings.Replace(uuid, "-", "", -1)
	return uuid, nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
