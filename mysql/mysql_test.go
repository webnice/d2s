package mysql

import (
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	d2sTypes "github.com/webnice/d2s/types"
)

func TestLog(t *testing.T) {
	if log() == nil {
		t.Fatal("log() вернул nil")
	}
}

func TestDialect(t *testing.T) {
	if Dialect() == nil {
		t.Fatal("Dialect() вернул nil")
	}
}

func TestDatabaseCurrent(t *testing.T) {
	msl := Dialect().(*impl)
	// Успешное получение имени базы данных.
	db := fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"DATABASE()"}, [][]driver.Value{{"testdb"}}, nil
	})
	if ret, err := msl.DatabaseCurrent(db); err != nil || ret != "testdb" {
		t.Fatalf("DatabaseCurrent() = %q, %v; ожидается testdb, nil", ret, err)
	}
	// Ошибка запроса.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return nil, nil, errors.New("сбой соединения")
	})
	if _, err := msl.DatabaseCurrent(db); err == nil || !strings.Contains(err.Error(), "ошибка запроса к базе данных") {
		t.Fatalf("ожидается ошибка запроса, получено: %v", err)
	}
	// Ошибка сканирования: две колонки вместо одной.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
	})
	if _, err := msl.DatabaseCurrent(db); err == nil || !strings.Contains(err.Error(), "ошибка сканирования") {
		t.Fatalf("ожидается ошибка сканирования, получено: %v", err)
	}
	// Пустой результат: база данных не выбрана.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"DATABASE()"}, nil, nil
	})
	if _, err := msl.DatabaseCurrent(db); err == nil || !strings.Contains(err.Error(), "база данных не выбрана") {
		t.Fatalf("ожидается ошибка 'база данных не выбрана', получено: %v", err)
	}
	// Ошибка закрытия строк попадает в лог, но не в результат.
	db = fakeDBCloseErr(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"DATABASE()"}, [][]driver.Value{{"testdb"}}, nil
	})
	if ret, err := msl.DatabaseCurrent(db); err != nil || ret != "testdb" {
		t.Fatalf("DatabaseCurrent() с ошибкой закрытия = %q, %v; ожидается testdb, nil", ret, err)
	}
	// Ошибка сканирования с прерыванием перебора: ошибка закрытия попадает в лог.
	db = fakeDBCloseErr(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
	})
	if _, err := msl.DatabaseCurrent(db); err == nil {
		t.Fatal("ожидается ошибка сканирования")
	}
}

func TestTableComment(t *testing.T) {
	msl := Dialect().(*impl)
	inf := &d2sTypes.TableInfo{Database: "testdb", Table: "t"}
	// Успешная загрузка комментария.
	db := fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "TABLE_COMMENT") {
			return []string{"TABLE_COMMENT"}, [][]driver.Value{{"Комментарий таблицы"}}, nil
		}
		return nil, nil, errors.New("неожиданный запрос")
	})
	if err := msl.TableComment(db, inf); err != nil || inf.Comment != "Комментарий таблицы" {
		t.Fatalf("TableComment() = %v, комментарий %q; ожидается nil, 'Комментарий таблицы'", err, inf.Comment)
	}
	// Ошибка запроса.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return nil, nil, errors.New("сбой")
	})
	if err := msl.TableComment(db, inf); err == nil || !strings.Contains(err.Error(), "ошибка запроса к базе данных") {
		t.Fatalf("ожидается ошибка запроса, получено: %v", err)
	}
	// Ошибка сканирования.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
	})
	if err := msl.TableComment(db, inf); err == nil || !strings.Contains(err.Error(), "ошибка сканирования") {
		t.Fatalf("ожидается ошибка сканирования, получено: %v", err)
	}
	// Ошибка закрытия строк попадает в лог.
	db = fakeDBCloseErr(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"TABLE_COMMENT"}, [][]driver.Value{{"c"}}, nil
	})
	if err := msl.TableComment(db, inf); err != nil {
		t.Fatalf("TableComment() с ошибкой закрытия = %v; ожидается nil", err)
	}
	// Ошибка сканирования с прерыванием перебора: ошибка закрытия попадает в лог.
	db = fakeDBCloseErr(t, func(query string) ([]string, [][]driver.Value, error) {
		return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
	})
	if err := msl.TableComment(db, inf); err == nil {
		t.Fatal("ожидается ошибка сканирования")
	}
}

// columnsRow Строка INFORMATION_SCHEMA.COLUMNS для скриптованного драйвера.
func columnsRow(name string, def driver.Value, nullable, dataType, columnType, key, extra, comment string,
	size, precision, scale driver.Value) []driver.Value {
	return []driver.Value{name, def, nullable, dataType, columnType, key, extra, comment, size, precision, scale}
}

func TestTableInfo(t *testing.T) {
	msl := Dialect().(*impl)
	script := func(t *testing.T, rows [][]driver.Value) func(string) ([]string, [][]driver.Value, error) {
		return func(query string) ([]string, [][]driver.Value, error) {
			switch {
			case strings.HasPrefix(query, "SELECT DATABASE()"):
				return []string{"DATABASE()"}, [][]driver.Value{{"testdb"}}, nil
			case strings.Contains(query, "TABLE_COMMENT"):
				return []string{"TABLE_COMMENT"}, [][]driver.Value{{"Таблица"}}, nil
			case strings.Contains(query, "FROM INFORMATION_SCHEMA.COLUMNS"):
				cols := []string{"COLUMN_NAME", "COLUMN_DEFAULT", "IS_NULLABLE", "DATA_TYPE", "COLUMN_TYPE",
					"COLUMN_KEY", "EXTRA", "COLUMN_COMMENT", "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION", "NUMERIC_SCALE"}
				return cols, rows, nil
			}
			return nil, nil, errors.New("неожиданный запрос: " + query)
		}
	}

	// Полная загрузка: имя базы данных из аргумента.
	db := fakeDB(t, script(t, [][]driver.Value{
		columnsRow("id", nil, "NO", "bigint", "bigint(20) unsigned", "PRI", "auto_increment", "",
			nil, int64(20), int64(0)),
		columnsRow("name", "аноним", "YES", "varchar", "varchar(255)", "", "",
			"Имя {nix.String}", int64(255), nil, nil),
	}))
	inf, err := msl.TableInfo(db, "testdb", "t")
	if err != nil {
		t.Fatalf("TableInfo() ошибка: %v", err)
	}
	if inf.Database != "testdb" || inf.Table != "t" || inf.Comment != "Таблица" || len(inf.Columns) != 2 {
		t.Fatalf("TableInfo() = %#v", inf)
	}
	id := inf.Columns[0]
	if id.Name != "id" || id.Default != nil || id.IsNullable || id.TypeSimple != "bigint" ||
		id.TypeFull != "BIGINT(20)" || !id.IsUnsigned || !id.IsPrimary || !id.IsAutoIncrement ||
		id.TypeGolang.Simple != d2sTypes.TUint64 || id.TypeGolang.String() != "uint64" ||
		id.Precision == nil || *id.Precision != 20 || id.Scale == nil || *id.Scale != 0 ||
		id.Size != nil {
		t.Fatalf("колонка id = %#v", id)
	}
	nm := inf.Columns[1]
	if nm.Name != "name" || nm.Default == nil || *nm.Default != "аноним" || !nm.IsNullable ||
		nm.Comment != "Имя {nix.String}" || nm.TypeGolang.Simple != d2sTypes.TString ||
		nm.TypeGolang.String() != "nix.String" || nm.Size == nil || *nm.Size != 255 {
		t.Fatalf("колонка name = %#v", nm)
	}

	// Пустое имя базы данных: определяется из соединения.
	db = fakeDB(t, script(t, nil))
	inf, err = msl.TableInfo(db, "", "t")
	if err != nil || inf.Database != "testdb" {
		t.Fatalf("TableInfo(пустая база) = %#v, %v", inf, err)
	}

	// Ошибка DatabaseCurrent.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.HasPrefix(query, "SELECT DATABASE()") {
			return nil, nil, errors.New("сбой")
		}
		return nil, nil, errors.New("неожиданный запрос")
	})
	if _, err = msl.TableInfo(db, "", "t"); err == nil {
		t.Fatal("ожидается ошибка DatabaseCurrent")
	}

	// Ошибка TableComment.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "TABLE_COMMENT") {
			return nil, nil, errors.New("сбой")
		}
		return []string{"DATABASE()"}, [][]driver.Value{{"testdb"}}, nil
	})
	if _, err = msl.TableInfo(db, "testdb", "t"); err == nil {
		t.Fatal("ожидается ошибка TableComment")
	}

	// Ошибка запроса колонок.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "FROM INFORMATION_SCHEMA.COLUMNS") {
			return nil, nil, errors.New("сбой")
		}
		return []string{"TABLE_COMMENT"}, [][]driver.Value{{"c"}}, nil
	})
	if _, err = msl.TableInfo(db, "testdb", "t"); err == nil ||
		!strings.Contains(err.Error(), "ошибка запроса к базе данных") {
		t.Fatalf("ожидается ошибка запроса колонок, получено: %v", err)
	}

	// Ошибка сканирования колонок: две колонки вместо одиннадцати.
	db = fakeDB(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "FROM INFORMATION_SCHEMA.COLUMNS") {
			return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
		}
		return []string{"TABLE_COMMENT"}, [][]driver.Value{{"c"}}, nil
	})
	if _, err = msl.TableInfo(db, "testdb", "t"); err == nil ||
		!strings.Contains(err.Error(), "ошибка сканирования") {
		t.Fatalf("ожидается ошибка сканирования колонок, получено: %v", err)
	}

	// Ошибка закрытия строк колонок попадает в лог.
	db = fakeDBCloseErr(t, script(t, [][]driver.Value{
		columnsRow("id", nil, "NO", "int", "int(11)", "PRI", "", "", nil, int64(10), int64(0)),
	}))
	if _, err = msl.TableInfo(db, "testdb", "t"); err != nil {
		t.Fatalf("TableInfo() с ошибкой закрытия = %v; ожидается nil", err)
	}

	// Ошибка сканирования колонок с прерыванием перебора: ошибка закрытия попадает в лог.
	db = fakeDBCloseErr(t, func(query string) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "FROM INFORMATION_SCHEMA.COLUMNS") {
			return []string{"a", "b"}, [][]driver.Value{{"x", "y"}}, nil
		}
		return []string{"TABLE_COMMENT"}, [][]driver.Value{{"c"}}, nil
	})
	if _, err = msl.TableInfo(db, "testdb", "t"); err == nil {
		t.Fatal("ожидается ошибка сканирования колонок")
	}
}

func TestTableColumnsFullTypeParse(t *testing.T) {
	msl := Dialect().(*impl)
	for _, tc := range []struct {
		full, wantType string
		wantUnsigned   bool
	}{
		{"bigint(20) unsigned", "BIGINT(20)", true},
		{"int(10) UNSIGNED", "INT(10)", true},
		{"double unsigned", "DOUBLE", true},
		{"int(11)", "INT(11)", false},
		{"varchar(255)", "VARCHAR(255)", false},
		{"text", "TEXT", false},
	} {
		col := &d2sTypes.ColumnInfo{TypeFull: tc.full}
		msl.tableColumnsFullTypeParse(col)
		if col.TypeFull != tc.wantType || col.IsUnsigned != tc.wantUnsigned {
			t.Fatalf("tableColumnsFullTypeParse(%q) = %q, %v; ожидается %q, %v",
				tc.full, col.TypeFull, col.IsUnsigned, tc.wantType, tc.wantUnsigned)
		}
	}
}

func TestColumnTypeWidth(t *testing.T) {
	msl := Dialect().(*impl)
	for _, tc := range []struct {
		full string
		want int
	}{
		{"TINYINT(1)", 1},
		{"TINYINT(4)", 4},
		{"BIT(8)", 8},
		{"TINYINT", 0},
		{"DECIMAL(10,2)", 0},
		{"", 0},
		{"INT(11", 0},
	} {
		col := &d2sTypes.ColumnInfo{TypeFull: tc.full}
		if got := msl.columnTypeWidth(col); got != tc.want {
			t.Fatalf("columnTypeWidth(%q) = %d; ожидается %d", tc.full, got, tc.want)
		}
	}
}

func TestColumnTypeMapping(t *testing.T) {
	msl := Dialect().(*impl)
	for _, tc := range []struct {
		name       string
		typeSimple string
		typeFull   string
		isUnsigned bool
		scale      *uint64
		want       string
	}{
		{"boolean", "boolean", "boolean(1)", false, nil, d2sTypes.TBool},
		{"bool", "bool", "bool", false, nil, d2sTypes.TBool},
		{"tinyint(1)", "tinyint", "tinyint(1)", false, nil, d2sTypes.TBool},
		{"tinyint(4)", "tinyint", "tinyint(4)", false, nil, d2sTypes.TInt8},
		{"tinyint без ширины", "tinyint", "tinyint", false, nil, d2sTypes.TInt8},
		{"binary", "binary", "binary(16)", false, nil, d2sTypes.TBytes},
		{"tinyblob", "tinyblob", "tinyblob", false, nil, d2sTypes.TBytes},
		{"blob", "blob", "blob", false, nil, d2sTypes.TBytes},
		{"mediumblob", "mediumblob", "mediumblob", false, nil, d2sTypes.TBytes},
		{"longblob", "longblob", "longblob", false, nil, d2sTypes.TBytes},
		{"varbinary", "varbinary", "varbinary(64)", false, nil, d2sTypes.TBytes},
		{"geometry", "geometry", "geometry", false, nil, d2sTypes.TBytes},
		{"point", "point", "point", false, nil, d2sTypes.TBytes},
		{"linestring", "linestring", "linestring", false, nil, d2sTypes.TBytes},
		{"polygon", "polygon", "polygon", false, nil, d2sTypes.TBytes},
		{"multipoint", "multipoint", "multipoint", false, nil, d2sTypes.TBytes},
		{"multilinestring", "multilinestring", "multilinestring", false, nil, d2sTypes.TBytes},
		{"multipolygon", "multipolygon", "multipolygon", false, nil, d2sTypes.TBytes},
		{"geometrycollection", "geometrycollection", "geometrycollection", false, nil, d2sTypes.TBytes},
		{"decimal(p,0)", "decimal", "decimal(20,0)", false, u64(0), d2sTypes.TBigInt},
		{"decimal без scale", "decimal", "decimal(20)", false, nil, d2sTypes.TBigInt},
		{"decimal(p,s)", "decimal", "decimal(10,2)", false, u64(2), d2sTypes.TFloat64},
		{"double", "double", "double", false, nil, d2sTypes.TFloat64},
		{"float", "float", "float", false, nil, d2sTypes.TFloat32},
		{"int", "int", "int(11)", false, nil, d2sTypes.TInt32},
		{"smallint", "smallint", "smallint(6)", false, nil, d2sTypes.TInt16},
		{"mediumint", "mediumint", "mediumint(9)", false, nil, d2sTypes.TInt32},
		{"bigint", "bigint", "bigint(20)", false, nil, d2sTypes.TInt64},
		{"int unsigned", "int", "int(10) unsigned", true, nil, d2sTypes.TUint32},
		{"smallint unsigned", "smallint", "smallint(5) unsigned", true, nil, d2sTypes.TUint16},
		{"mediumint unsigned", "mediumint", "mediumint(8) unsigned", true, nil, d2sTypes.TUint32},
		{"bigint unsigned", "bigint", "bigint(20) unsigned", true, nil, d2sTypes.TUint64},
		{"bit(1)", "bit", "bit(1)", false, nil, d2sTypes.TBool},
		{"bit без ширины", "bit", "bit", false, nil, d2sTypes.TBool},
		{"bit(8)", "bit", "bit(8)", false, nil, d2sTypes.TBytes},
		{"year", "year", "year(4)", false, nil, d2sTypes.TInt16},
		{"char", "char", "char(10)", false, nil, d2sTypes.TString},
		{"enum", "enum", "enum('a','b')", false, nil, d2sTypes.TString},
		{"varchar", "varchar", "varchar(255)", false, nil, d2sTypes.TString},
		{"tinytext", "tinytext", "tinytext", false, nil, d2sTypes.TString},
		{"text", "text", "text", false, nil, d2sTypes.TString},
		{"mediumtext", "mediumtext", "mediumtext", false, nil, d2sTypes.TString},
		{"longtext", "longtext", "longtext", false, nil, d2sTypes.TString},
		{"set", "set", "set('x','y')", false, nil, d2sTypes.TStringSlice},
		{"json", "json", "json", false, nil, d2sTypes.TRawMessage},
		{"date", "date", "date", false, nil, d2sTypes.TTime},
		{"datetime", "datetime", "datetime", false, nil, d2sTypes.TTime},
		{"timestamp", "timestamp", "timestamp", false, nil, d2sTypes.TTime},
		{"time", "time", "time", false, nil, d2sTypes.TDuration},
		{"uuid", "uuid", "uuid", false, nil, d2sTypes.TUUID},
		{"inet4", "inet4", "inet4", false, nil, d2sTypes.TIP},
		{"inet6", "inet6", "inet6", false, nil, d2sTypes.TIP},
		{"vector", "vector", "vector(3)", false, nil, d2sTypes.TBytes},
		{"неизвестный тип", "exotic_type", "exotic", false, nil, d2sTypes.TUnknown},
	} {
		col := &d2sTypes.ColumnInfo{
			TypeSimple: tc.typeSimple,
			TypeFull:   tc.typeFull,
			IsUnsigned: tc.isUnsigned,
			Scale:      tc.scale,
		}
		msl.columnTypeMapping(col)
		if col.TypeGolang.Simple != tc.want {
			t.Fatalf("%s: Simple = %q; ожидается %q", tc.name, col.TypeGolang.Simple, tc.want)
		}
	}
	// Nullable вариант типа.
	col := &d2sTypes.ColumnInfo{TypeSimple: "int", TypeFull: "INT(11)", IsNullable: true}
	msl.columnTypeMapping(col)
	if got := col.TypeGolang.String(); got != "nix.Int32" {
		t.Fatalf("nullable int = %q; ожидается nix.Int32", got)
	}
}

func u64(v uint64) *uint64 { return &v }
