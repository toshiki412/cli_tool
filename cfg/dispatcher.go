package cfg

import (
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func DispatchTarget(target TargetType, table TargetFuncTable) {
	switch target.Kind {
	case "mysql":
		var conf TargetMysqlType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		defaults.SetDefaults(&conf)
		table.Mysql(conf)
	case "file":
		var conf TargetFileType
		err := mapstructure.Decode(target.Config, &conf)
		cobra.CheckErr(err)
		defaults.SetDefaults(&conf)
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
		defaults.SetDefaults(&conf)
		table.Gcs(conf)
	default:
		panic("unknown storage kind")
	}
}
