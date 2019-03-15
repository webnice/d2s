package d2s // import "gopkg.in/webnice/d2s.v1/d2s"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
)

// New creates a new object and return interface
func New() Interface {
	var d2s = new(impl)
	return d2s
}

// Debug sets a debug mode
func (d2s *impl) Debug(d bool) Interface { d2s.debug = d; return d2s }

// Dialect sets the SQL dialect
func (d2s *impl) Dialect(sqlDialect string) (err error) {
	switch sqlDialect {
	case DrvSqlite3:
		d2s.dialect = &Sqlite3Dialect{}
	case DrvMysql:
		d2s.dialect = &MySQLDialect{}
	case DrvClickhouse:
		d2s.dialect = &ClickHouseDialect{}
	case DrvRedshift:
		d2s.dialect = &RedshiftDialect{}
	case DrvTidb:
		d2s.dialect = &TiDBDialect{}
	case DrvPostgres:
		d2s.dialect = &PostgresDialect{}
	default:
		err = fmt.Errorf("unknown dialect: %q", sqlDialect)
		return
	}

	return
}

// Run runs a command
func (d2s *impl) Run(command string, tableName string, structureName string, fileName string) (err error) {

	return
}
