package cfg

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func DispatchTarget(target TargetType, table TargetFuncTable) {
	switch target.Kind {
	case "mysql":
		var conf TargetMysqlType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		table.Mysql(conf)
	case "file":
		var conf TargetFileType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		table.File(conf)
	default:
		panic("unknown target kind")
	}
}

func DispatchStorages(storage StorageType, table StorageFuncTable) {
	switch storage.Kind {
	case "gcs":
		var conf StorageGoogleStorageType
		err := mapstructure.Decode(storage.Config, &conf)
		cobra.CheckErr(err)
		table.Gcs(conf)
	default:
		panic("unknown storage kind")
	}
}
