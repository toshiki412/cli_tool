package cfg // フォルダ名

type TragetMysqlConfigType struct {
	Host string
	Port int
	User string
	Password string
	Database string
}

type TargetType struct {
	Kind string
	Config interface{}
}

type UploadType struct {
	Kind string
	Config interface{}
}

type ConfigType struct {
	Target TargetType // 頭文字を大文字にすることで、外部からアクセス可能になる. 小文字だとプライベートになる
	Upload UploadType
}