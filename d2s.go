package d2s

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/webnice/d2s/mysql"
	d2sTypes "github.com/webnice/d2s/types"
)

// New Создание нового объекта и возврат интерфейса.
func New() Interface {
	var d2s = new(impl)
	return d2s
}

// Debug Установка режима отладки.
func (d2s *impl) Debug(d bool) Interface { d2s.debug = d; return d2s }

// Types Переопределение типов полей структуры.
// Список "поле1:тип;поле2:тип", где поле — имя колонки базы данных (без учёта
// регистра), тип — стандартный тип golang или тип nix/v2. Переопределение
// имеет наивысший приоритет над переопределением из комментария колонки
// и базовым сопоставлением типов.
func (d2s *impl) Types(overrideList string) (err error) {
	if strings.TrimSpace(overrideList) == "" {
		return
	}
	if d2s.override, err = d2sTypes.ParseOverrideList(overrideList); err != nil {
		return
	}

	return
}

// Dialect Установка диалекта SQL.
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
		err = fmt.Errorf("не реализован диалект: %q", sqlDialect)
		return
	}

	return
}

// Connect Соединение с базой данных для получения информации.
func (d2s *impl) Connect(db *sql.DB) Interface { d2s.db = db; return d2s }

// Create Создание структуры по таблице базы данных.
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
		err = fmt.Errorf("ошибка получения информации о таблице: %s", err)
		return
	}
	inf.Package, inf.Struct = packageName, structureName
	if buf, err = d2s.Generator(inf); err != nil {
		err = fmt.Errorf("ошибка генерации структуры: %s", err)
		return
	}
	if fh, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666)); err != nil {
		err = fmt.Errorf("ошибка открытия файла %q: %s", fileName, err)
		return
	}
	defer func() {
		if e := fh.Close(); e != nil {
			err = e
		}
	}()
	if _, err = buf.WriteTo(fh); err != nil {
		err = fmt.Errorf("ошибка записи в файл %q: %s", fileName, err)
		return
	}

	return
}
