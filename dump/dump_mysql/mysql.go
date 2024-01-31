package dump_mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/format"
	_ "github.com/pingcap/tidb/parser/test_driver"

	"github.com/JamesStewy/go-mysqldump"

	"github.com/briandowns/spinner"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/toshiki412/cli_tool/cfg"
)

func MysqlDumpFile(dumpDir string, conf cfg.TargetMysqlType) string {
	dumpFilename := fmt.Sprintf("%s-%s.sql", "mysql", conf.Database)
	dumpFile := filepath.Join(dumpDir, dumpFilename)
	return dumpFile
}

func Dump(dumpFile string, conf cfg.TargetMysqlType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" mysql dump ... (database: %s)", conf.Database)
	s.Start()

	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.DBName = conf.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	filename := filepath.Base(dumpFile)
	dumpDir := filepath.Dir(dumpFile)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	dumper, err := mysqldump.Register(db, dumpDir, filename)
	cobra.CheckErr(err)

	_, err = dumper.Dump()
	cobra.CheckErr(err)

	dumper.Close()

	s.Stop()
	fmt.Printf("✔︎ mysql dump completed. (database: %s)\n", conf.Database)
}

func Import(dumpFile string, conf cfg.TargetMysqlType) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" mysql import ... (database: %s)", conf.Database)
	s.Start()

	config := mysql.NewConfig()
	config.User = conf.User
	config.Passwd = conf.Password
	config.DBName = conf.Database
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	content, err := os.ReadFile(dumpFile)
	cobra.CheckErr(err)

	p := parser.New()

	stmts, _, err := p.Parse(string(content), "", "")
	if err != nil {
		log.Fatalf("failed to parse seed sql: %v", err)
	}

	var buf bytes.Buffer
	for _, stmt := range stmts {
		buf.Reset()

		// 各ast.StmtNodeをSQL文字列に復元する
		stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf))

		sql := buf.String()
		if _, err := db.Exec(sql); err != nil {
			log.Fatalf("failed to execute sql: err=%v sql=%s", err, sql[:100])
		}
	}
	s.Stop()
	fmt.Printf("✔︎ mysql import completed. (database: %s)\n", conf.Database)
}
