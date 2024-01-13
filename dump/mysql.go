package dump

import (
	"database/sql"
	"fmt"

	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
)

func Dump(dumpDir string, conf cfg.TargetMysqlType) {
	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.DBName = conf.Database

	dumpFileNameFormat := fmt.Sprintf("%s-%s", "mysql", conf.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFileNameFormat, config.DBName)
	cobra.CheckErr(err)

	// dump database to file
	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Println("successfully dumped to file.")

	dumper.Close()
}
