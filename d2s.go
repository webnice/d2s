// Package d2s
package d2s

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"

	"github.com/webnice/d2s/mysql"
	d2sTypes "github.com/webnice/d2s/types"
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
	case d2sTypes.DrvMysql:
		d2s.dialect = mysql.Dialect()
	//case DrvSqlite3:
	//	d2s.dialect = &Sqlite3Dialect{}
	//case DrvClickhouse:
	//	d2s.dialect = &ClickHouseDialect{}
	//case DrvRedshift:
	//	d2s.dialect = &RedshiftDialect{}
	//case DrvTidb:
	//	d2s.dialect = &TiDBDialect{}
	//case DrvPostgres:
	//	d2s.dialect = &PostgresDialect{}
	default:
		err = fmt.Errorf("not implemented dialect: %q", sqlDialect)
		return
	}

	return
}

// Connect Соединение с базой данных для получения информации
func (d2s *impl) Connect(db *sql.DB) Interface { d2s.db = db; return d2s }

// Create structure from table
func (d2s *impl) Create(
	databaseName string,
	tableName string,
	packageName string,
	structureName string,
	fileName string,
) (err error) {
	var (
		inf *d2sTypes.TableInfo
		buf *bytes.Buffer
		fh  *os.File
	)

	if inf, err = d2s.dialect.TableInfo(d2s.db, databaseName, tableName); err != nil {
		err = fmt.Errorf("get table info error: %s", err)
		return
	}
	inf.Package, inf.Struct = packageName, structureName
	if buf, err = d2s.Generator(inf); err != nil {
		err = fmt.Errorf("generate structure error: %s", err)
		return
	}
	if fh, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666)); err != nil {
		err = fmt.Errorf("open file %q error: %s", fileName, err)
		return
	}
	defer func() {
		if e := fh.Close(); e != nil {
			err = e
		}
	}()
	if _, err = buf.WriteTo(fh); err != nil {
		err = fmt.Errorf("write to file %q error: %s", fileName, err)
		return
	}

	return
}
