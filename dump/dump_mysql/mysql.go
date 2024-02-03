package dump_mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/format"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"

	"github.com/aliakseiz/go-mysqldump"

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

	filename := strings.Replace(filepath.Base(dumpFile), ".sql", "", 1)
	dumpDir := filepath.Dir(dumpFile)

	db, err := sql.Open("mysql", config.FormatDSN())
	cobra.CheckErr(err)

	dumper, err := mysqldump.Register(db, dumpDir, filename, conf.Database)
	cobra.CheckErr(err)

	err = dumper.Dump()
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
		fmt.Printf("failed to parse seed sql: %v\n", err)
	}

	var buf bytes.Buffer
	for _, stmt := range stmts {
		buf.Reset()

		// 各ast.StmtNodeをSQL文字列に復元する
		stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf))

		sql := buf.String()
		if _, err := db.Exec(sql); err != nil {
			fmt.Printf("failed to execute sql: err=%v sql=%s\n", err, sql[:100])
		}
	}
	s.Stop()
	fmt.Printf("✔︎ mysql import completed. (database: %s)\n", conf.Database)
}
