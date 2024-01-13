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

type UploadGoogleStorageType struct {
	Bucket string
	Dir    string
}

type UploadType struct {
	Kind   string
	Config interface{}
}

type SettingType struct {
	Target TargetType // 頭文字を大文字にすることで、外部からアクセス可能になる. 小文字だとプライベートになる
	Upload UploadType
}

type VersionType struct {
	Id      string `json:"id"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}
