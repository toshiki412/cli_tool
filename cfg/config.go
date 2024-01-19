package cfg // フォルダ名

import "path/filepath"

type TargetMysqlType struct {
	Host     string `default:"localhost"`
	Port     int    `default:"3306"`
	User     string `default:"root"`
	Password string `default:""`
	Database string
}

type TargetFileType struct {
	Path string
}

type TargetType struct {
	Kind   string
	Config interface{}
}

type StorageGoogleStorageType struct {
	Bucket string
	Dir    string
}

type StorageType struct {
	Kind   string
	Config interface{}
}

type SettingType struct {
	Targets []TargetType // 頭文字を大文字にすることで、外部からアクセス可能になる. 小文字だとプライベートになる
	Storage StorageType
}

type VersionType struct {
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

func (v VersionType) CreateZipFile() string {
	return v.Id + ".zip"
}

func (v VersionType) CreateZipFileWithDir(dir string) string {
	return filepath.Join(dir, v.CreateZipFile())
}

type TargetMysqlFunc func(config TargetMysqlType) // 関数型
type TargetFileFunc func(config TargetFileType)
type TargetFuncTable struct {
	Mysql TargetMysqlFunc
	File  TargetFileFunc
}

type StorageGoogleStorageFunc func(config StorageGoogleStorageType)
type StorageFuncTable struct {
	Gcs StorageGoogleStorageFunc
}

type DataType struct {
	Version   string        `json:"version"`
	Histories []VersionType `json:"histories"`
}
