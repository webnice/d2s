package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	d2sTypes "github.com/webnice/d2s/v2/types"
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
		err = fmt.Errorf("ошибка запроса к базе данных: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("ошибка закрытия строк запроса: %s", e)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&ret); err != nil {
			err = fmt.Errorf("ошибка сканирования: %s", err)
			return
		}
	}
	if ret == "" {
		err = fmt.Errorf("база данных не выбрана")
		return
	}

	return
}

// TableComment Загрузка комментариев к таблице.
func (msl *impl) TableComment(db *sql.DB, inf *d2sTypes.TableInfo) (err error) {
	const dbQuery = "SELECT TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	var rows *sql.Rows

	if rows, err = db.Query(dbQuery, inf.Database, inf.Table); err != nil {
		err = fmt.Errorf("ошибка запроса к базе данных: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("ошибка закрытия строк запроса: %s", e)
		}
	}()
	for rows.Next() {
		if err = rows.Scan(&inf.Comment); err != nil {
			err = fmt.Errorf("ошибка сканирования: %s", err)
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
		err = fmt.Errorf("ошибка запроса к базе данных: %s", err)
		return
	}
	defer func() {
		if e := rows.Close(); e != nil {
			log().Errorf("ошибка закрытия строк запроса: %s", e)
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
		// Признак nullable.
		isNullable = strings.TrimSpace(strings.ToLower(isNullable))
		if isNullable == boolTrue1 || isNullable == boolTrue2 {
			row.IsNullable = true
		}
		// Полный тип колонки и признак UNSIGNED.
		msl.tableColumnsFullTypeParse(row)
		// Первичный ключ.
		isPrimary = strings.TrimSpace(strings.ToLower(isPrimary))
		row.IsPrimary = isPrimary == keyPrimary
		// Автоинкремент.
		row.IsAutoIncrement = strings.Contains(strings.ToLower(isAutoIncrement), keyAutoIncrement)
		// Если что-то пошло не так.
		if err != nil {
			err = fmt.Errorf("ошибка сканирования: %s", err)
			return
		}
		// Конвертация типа в тип golang.
		msl.columnTypeMapping(row)
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
func (msl *impl) columnTypeMapping(col *d2sTypes.ColumnInfo) {
	switch strings.ToLower(col.TypeSimple) {
	case dbBoolean, dbBool, dbTinyint:
		// tinyint(1) — булев тип, tinyint без ширины и tinyint(N>1) — числовой.
		if strings.ToLower(col.TypeSimple) == dbTinyint && msl.columnTypeWidth(col) != 1 {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TInt8)
		} else {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBool)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbBinary, dbTinyblob, dbBlob, dbMediumblob, dbLongblob, dbVarbinary:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBytes)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbGeometry, dbPoint, dbLinestring, dbPolygon, dbMultipoint, dbMultilinestring,
		dbMultipolygon, dbGeometrycollection:
		// Пространственные типы отдаются драйвером во внутреннем бинарном формате.
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBytes)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbDecimal:
		// decimal(p,0) — целочисленный десятичный тип без потери точности.
		if col.Scale != nil && *col.Scale > 0 {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TFloat64)
		} else {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBigInt)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbDouble, dbFloat:
		// float — 32-битный, double — 64-битный.
		if strings.ToLower(col.TypeSimple) == dbDouble {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TFloat64)
		} else {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TFloat32)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbInt, dbSmallint, dbMediumint, dbBigint:
		var signed, unsigned = d2sTypes.TInt32, d2sTypes.TUint32
		switch strings.ToLower(col.TypeSimple) {
		case dbSmallint:
			signed, unsigned = d2sTypes.TInt16, d2sTypes.TUint16
		case dbBigint:
			signed, unsigned = d2sTypes.TInt64, d2sTypes.TUint64
		}
		if col.IsUnsigned {
			col.TypeGolang = d2sTypes.NewBase(unsigned)
		} else {
			col.TypeGolang = d2sTypes.NewBase(signed)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbBit:
		// bit(1) — булев тип, bit(N>1) — битовая строка.
		if msl.columnTypeWidth(col) > 1 {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBytes)
		} else {
			col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBool)
		}
		col.TypeGolang.Nullable(col.IsNullable)
	case dbYear:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TInt16)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbChar, dbEnum, dbVarchar, dbTinytext, dbText, dbMediumtext, dbLongtext:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TString)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbSet:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TStringSlice)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbJson:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TRawMessage)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbDate, dbDatetime, dbTimestamp:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TTime)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbTime:
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TDuration)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbUuid:
		// Нативный тип UUID MariaDB 10.7+.
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TUUID)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbInet4, dbInet6:
		// Типы адресов IPv4 и IPv6 MariaDB 10.5+.
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TIP)
		col.TypeGolang.Nullable(col.IsNullable)
	case dbVector:
		// Векторный тип MariaDB 11.7+.
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TBytes)
		col.TypeGolang.Nullable(col.IsNullable)
	default:
		// Неизвестный тип базы данных: тип поля можно переопределить
		// через командную строку или комментарий колонки.
		col.TypeGolang = d2sTypes.NewBase(d2sTypes.TUnknown)
		col.TypeGolang.Nullable(col.IsNullable)
	}
}

// columnTypeWidth Извлечение ширины типа из полного типа колонки.
// Например, из "TINYINT(1)" — 1, из "BIT(8)" — 8. Если ширина не указана — 0.
func (msl *impl) columnTypeWidth(col *d2sTypes.ColumnInfo) (width int) {
	const brackets = `()`

	var open, close int

	if open = strings.Index(col.TypeFull, brackets[:1]); open < 0 {
		return
	}
	if close = strings.Index(col.TypeFull, brackets[1:]); close <= open {
		return
	}
	width, _ = strconv.Atoi(col.TypeFull[open+1 : close])

	return
}
