package dump_mysql

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/aliakseiz/go-mysqldump"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/tanimutomo/sqlfile"
	"github.com/toshiki412/cli_tool/cfg"
)

func Dump(dumpDir string, conf cfg.TargetMysqlType) {
	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.DBName = conf.Database

	dumpFileName := fmt.Sprintf("%s-%s", "mysql", conf.Database)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	// register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFileName, config.DBName)
	cobra.CheckErr(err)

	// dump database to file
	err = dumper.Dump()
	cobra.CheckErr(err)

	fmt.Println("successfully dumped to file.")

	dumper.Close()
}

func Import(dumpDir string, conf cfg.TargetMysqlType) {
	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	config.DBName = conf.Database

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	dumpFileName := fmt.Sprintf("%s-%s.sql", "mysql", conf.Database)
	dumpFile := filepath.Join(dumpDir, dumpFileName)

	s := sqlfile.New()
	err = s.File(dumpFile)
	cobra.CheckErr(err)

	_, err = s.Exec(db)
	cobra.CheckErr(err)
}
