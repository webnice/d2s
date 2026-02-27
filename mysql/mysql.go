package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	d2sTypes "github.com/webnice/d2s/types"
	kitModuleCfg "github.com/webnice/kit/v4/module/cfg"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Ссылка на менеджер логирования.
func log() kitTypes.Logger { return kitModuleCfg.Get().Log() }

// Dialect Возвращает диалект SQL.
func Dialect() Interface {
	var msl = new(impl)
	return msl
}

// TableInfo Получение информации о структуре таблицы.
func (msl *impl) TableInfo(db *sql.DB, databaseName string, tableName string) (ret *d2sTypes.TableInfo, err error) {
	// Структура ответа.
	ret = &d2sTypes.TableInfo{
		Database: databaseName,
		Table:    tableName,
	}
	if ret.Database == "" {
		if ret.Database, err = msl.DatabaseCurrent(db); err != nil {
			return
		}
	}
	// Загрузка комментариев к таблице.
	if err = msl.TableComment(db, ret); err != nil {
		return
	}
	// Загрузка колонок таблицы.
	if err = msl.tableColumns(db, ret); err != nil {
		return
	}

	return
}

// DatabaseCurrent Получение имени текущей базы данных выбранной через DSN.
func (msl *impl) DatabaseCurrent(db *sql.DB) (ret string, err error) {
	const dbQuery = `SELECT DATABASE()`
	var rows *sql.Rows

	if rows, err = db.Query(dbQuery); err != nil {
		err = fmt.Errorf("query database error: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("close query rows error: %s", e)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&ret); err != nil {
			err = fmt.Errorf("scan error: %s", err)
			return
		}
	}
	if ret == "" {
		err = fmt.Errorf("database is not selected")
		return
	}

	return
}

// TableComment Загрузка комментариев к таблице.
func (msl *impl) TableComment(db *sql.DB, inf *d2sTypes.TableInfo) (err error) {
	const dbQuery = "SELECT TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	var rows *sql.Rows

	if rows, err = db.Query(dbQuery, inf.Database, inf.Table); err != nil {
		err = fmt.Errorf("query database error: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("close query rows error: %s", e)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&inf.Comment); err != nil {
			err = fmt.Errorf("scan error: %s", err)
			return
		}
	}

	return
}

// Загрузка колонок таблицы.
func (msl *impl) tableColumns(db *sql.DB, inf *d2sTypes.TableInfo) (err error) {
	const (
		dbQuery = "SELECT COLUMN_NAME, COLUMN_DEFAULT, IS_NULLABLE, DATA_TYPE, COLUMN_TYPE, COLUMN_KEY" +
			", EXTRA, COLUMN_COMMENT, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE" +
			" FROM INFORMATION_SCHEMA.COLUMNS" +
			" WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?" +
			" ORDER BY INFORMATION_SCHEMA.COLUMNS.ORDINAL_POSITION"
		boolTrue1        = `yes`
		boolTrue2        = `true`
		keyPrimary       = `pri`
		keyAutoIncrement = `auto_increment`
	)
	var (
		rows                                   *sql.Rows
		row                                    *d2sTypes.ColumnInfo
		isNullable, isPrimary, isAutoIncrement string
	)

	if rows, err = db.Query(dbQuery, inf.Database, inf.Table); err != nil {
		err = fmt.Errorf("query database error: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("close query rows error: %s", e)
		}
	}()
	for rows.Next() {
		row = new(d2sTypes.ColumnInfo)
		err = rows.Scan(
			&row.Name,
			&row.Default,
			&isNullable,
			&row.TypeSimple,
			&row.TypeFull,
			&isPrimary,
			&isAutoIncrement,
			&row.Comment,
			&row.Size,
			&row.Precision,
			&row.Scale,
		)
		// NULLABLE.
		isNullable = strings.TrimSpace(strings.ToLower(isNullable))
		if isNullable == boolTrue1 || isNullable == boolTrue2 {
			row.IsNullable = true
		}
		// Полный тип и UNSIGNED расширение типа.
		msl.tableColumnsFullTypeParse(row)
		// Первичный ключ.
		isPrimary = strings.TrimSpace(strings.ToLower(isPrimary))
		row.IsPrimary = isPrimary == keyPrimary
		// Авто инкремент.
		row.IsAutoIncrement = strings.Contains(strings.ToLower(isAutoIncrement), keyAutoIncrement)
		// Если что-то пошло не так
		if err != nil {
			err = fmt.Errorf("scan error: %s", err)
			return
		}
		// Конвертация типа в тип Golang.
		if err = msl.columnTypeMapping(row); err != nil {
			err = fmt.Errorf("database type conversion to golang type error: %s", err)
		}
		inf.Columns = append(inf.Columns, row)
	}

	return
}

func (msl *impl) tableColumnsFullTypeParse(col *d2sTypes.ColumnInfo) {
	const keyUnsigned = `unsigned`
	var tmp []string

	if tmp = strings.Split(col.TypeFull, " "); len(tmp) == 2 {
		col.TypeFull = tmp[0]
		col.IsUnsigned = strings.TrimSpace(strings.ToLower(tmp[1])) == keyUnsigned
	}
	col.TypeFull = strings.ToUpper(col.TypeFull)
}

// Сопоставление типов данных в БД с типом данных в golang.
func (msl *impl) columnTypeMapping(col *d2sTypes.ColumnInfo) (err error) {
	switch strings.ToLower(col.TypeSimple) {
	case dbBoolean, dbBool, dbTinyint:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBool)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbBinary, dbTinyblob, dbBlob, dbMediumblob, dbLongblob, dbVarbinary:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBytes)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbDecimal, dbDouble, dbFloat:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TFloat64)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbInt, dbSmallint, dbMediumint, dbBigint:
		if col.IsUnsigned {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TUint64)
		} else {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TInt64)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbChar, dbEnum, dbVarchar, dbTinytext, dbText, dbMediumtext, dbLongtext:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TString)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbDate, dbDatetime, dbTime, dbTimestamp:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TTime)
		col.TypeGolang.Nullable(col.IsNullable)
	default:
		err = fmt.Errorf("not implemented type mapping for type %q", col.TypeSimple)
		return
	}

	return
}
