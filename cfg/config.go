package cfg // フォルダ名

type TargetMysqlType struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
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
	Target  TargetType // 頭文字を大文字にすることで、外部からアクセス可能になる. 小文字だとプライベートになる
	Storage StorageType
}

type VersionType struct {
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

type TargetMysqlFunc func(config TargetMysqlType) // 関数型
type TargetFuncTable struct {
	Mysql TargetMysqlFunc
}

type StorageGoogleStorageFunc func(config StorageGoogleStorageType)
type StorageFuncTable struct {
	Gcs StorageGoogleStorageFunc
}
