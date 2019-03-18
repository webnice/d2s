package main

//import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"database/sql"
	"fmt"

	"gopkg.in/webnice/d2s.v1/d2s"
	d2sTypes "gopkg.in/webnice/d2s.v1/d2s/types"

	// Init database drivers
	_ "github.com/go-sql-driver/mysql" // Mysql
	//_ "github.com/mattn/go-sqlite3"     // Sqlite
	//_ "github.com/kshvakov/clickhouse"  // Clickhouse
	//_ "github.com/lib/pq"               // Postgres, Cockroach, Redshift
	//_ "github.com/ziutek/mymysql/godrv" // App Engine CloudSQL
)

func main() {
	const cmdCreate = `create`
	var err error
	var db2struct d2s.Interface
	var cmd string
	var arg *Args
	var db *sql.DB

	// Логирование
	log.Gist().StandardLogSet()
	defer log.Done()
	defer log.Gist().StandardLogUnset()
	// Checking driver and set dialect
	db2struct = d2s.New()
	switch cmd, arg = args(); arg.Driver {
	case d2sTypes.DrvMysql, d2sTypes.DrvPostgres, d2sTypes.DrvSqlite3, d2sTypes.DrvClickhouse:
		err = db2struct.Dialect(arg.Driver)
	case d2sTypes.DrvRedshift:
		err = db2struct.Dialect(arg.Driver)
		arg.Driver = d2sTypes.DrvPostgres
	case d2sTypes.DrvTidb:
		err = db2struct.Dialect(arg.Driver)
		arg.Driver = d2sTypes.DrvMysql
	default:
		err = fmt.Errorf("%q driver not supported", arg.Driver)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	// Opening database connection
	if db, err = sql.Open(arg.Driver, arg.Dsn); err != nil {
		log.Fatalf("connect to database error: %s", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("close database error: %s", err)
		}
	}()
	db.SetConnMaxLifetime(0)
	// Check correct value
	if arg.Table == "" || arg.Package == "" || arg.Structure == "" || arg.File == "" {
		argUsage()
		return
	}
	// Running command of database migration with different arguments
	switch cmd {
	case cmdCreate:
		err = db2struct.
			Connect(db).
			Create(arg.Database, arg.Table, arg.Package, arg.Structure, arg.File)
	default:
		argUsage()
	}
	if err != nil {
		log.Fatalf("db2struct error: %s", err)
	}
}
