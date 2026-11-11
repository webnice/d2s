package main

import (
	"database/sql"
	"fmt"

	"github.com/webnice/d2s/v2"
	d2sTypes "github.com/webnice/d2s/v2/types"
	kitModuleCfg "github.com/webnice/kit/v4/module/cfg"
	kitTypes "github.com/webnice/kit/v4/types"

	// Драйвера базы данных.
	_ "github.com/go-sql-driver/mysql" // Mysql.
	//_ "github.com/mattn/go-sqlite3"     // Sqlite.
	//_ "github.com/kshvakov/clickhouse"  // Clickhouse.
	//_ "github.com/lib/pq"               // Postgres, Cockroach, Redshift.
	//_ "github.com/ziutek/mymysql/godrv" // App Engine CloudSQL.
)

// Ссылка на менеджер логирования.
func log() kitTypes.Logger { return kitModuleCfg.Get().Log() }

func main() {
	const cmdCreate = `create`
	var (
		err       error
		db2struct d2s.Interface
		cmd       string
		arg       *Args
		db        *sql.DB
	)

	// Проверка драйвера и установка диалекта.
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
		err = fmt.Errorf("драйвер %q не поддерживается", arg.Driver)
	}
	if err != nil {
		log().Fatal(err.Error())

		return
	}
	// Открытие соединения с базой данных.
	if db, err = sql.Open(arg.Driver, arg.Dsn); err != nil {
		log().Fatalf("ошибка соединения с базой данных: %s", err)

		return
	}
	defer func() {
		if err = db.Close(); err != nil {
			log().Fatalf("ошибка закрытия соединения с базой данных: %s", err)
		}
	}()
	db.SetConnMaxLifetime(0)
	// Проверка корректности значений.
	if arg.Table == "" || arg.Package == "" || arg.Structure == "" || arg.File == "" {
		argUsage()
		return
	}
	// Выполнение команды миграции базы данных с различными аргументами.
	switch cmd {
	case cmdCreate:
		if err = db2struct.Types(arg.Types); err != nil {
			log().Fatalf("ошибка переопределения типов: %s", err)

			return
		}
		err = db2struct.
			Connect(db).
			Create(arg.Database, arg.Table, arg.Package, arg.Structure, arg.File)
	default:
		argUsage()
	}
	if err != nil {
		log().Fatalf("ошибка db2struct: %s", err)
	}
}
