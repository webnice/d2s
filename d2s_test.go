package d2s

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	d2sTypes "github.com/webnice/d2s/types"
)

// mockDialect Реализация интерфейса d2sTypes.Dialect для тестов.
type mockDialect struct {
	inf *d2sTypes.TableInfo
	err error
}

func (m *mockDialect) DatabaseCurrent(db *sql.DB) (string, error) { return "testdb", m.err }
func (m *mockDialect) TableInfo(db *sql.DB, databaseName, tableName string) (*d2sTypes.TableInfo, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.inf, nil
}

func TestNewDebugConnect(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("New() вернул nil")
	}
	// Цепочка вызовов возвращает тот же объект.
	if d.Debug(true) != d {
		t.Fatal("Debug() вернул другой объект")
	}
	if d.Connect(nil) != d {
		t.Fatal("Connect() вернул другой объект")
	}
	// Debug меняет флаг.
	impl := d.(*impl)
	d.Debug(true)
	if !impl.debug {
		t.Fatal("Debug(true) не установил флаг")
	}
}

func TestTypes(t *testing.T) {
	d := New()
	// Пустой список — без ошибки и без переопределений.
	if err := d.Types(""); err != nil {
		t.Fatalf("Types(пустой) = %v; ожидается nil", err)
	}
	// Корректный список.
	if err := d.Types("id:uint64;payload:nix.RawMessage"); err != nil {
		t.Fatalf("Types(корректный) = %v; ожидается nil", err)
	}
	impl := d.(*impl)
	if impl.override["id"] != "uint64" || impl.override["payload"] != "nix.RawMessage" {
		t.Fatalf("override = %#v", impl.override)
	}
	// Неизвестный тип — ошибка.
	if err := d.Types("id:uint128"); err == nil {
		t.Fatal("Types(неизвестный тип) = nil; ожидается ошибка")
	}
}

func TestDialect(t *testing.T) {
	d := New()
	if err := d.Dialect(d2sTypes.DrvMysql); err != nil {
		t.Fatalf("Dialect(mysql) = %v; ожидается nil", err)
	}
	if err := d.Dialect("oracle"); err == nil || !strings.Contains(err.Error(), "не реализован диалект") {
		t.Fatalf("Dialect(oracle) = %v; ожидается ошибка 'не реализован диалект'", err)
	}
}

func TestCreate(t *testing.T) {
	newImpl := func(m *mockDialect) *impl {
		d := New().(*impl)
		d.dialect = m

		return d
	}
	tableInfo := func() *d2sTypes.TableInfo {
		return &d2sTypes.TableInfo{
			Database: "testdb", Table: "message", Comment: "Сообщение",
			Columns: []*d2sTypes.ColumnInfo{
				{Name: "id", TypeGolang: d2sTypes.NewBase(d2sTypes.TUint64), TypeFull: "BIGINT(20)",
					IsPrimary: true, IsAutoIncrement: true},
			},
		}
	}

	// Успешное создание файла.
	fileName := filepath.Join(t.TempDir(), "message_model.go")
	d := newImpl(&mockDialect{inf: tableInfo()})
	if err := d.Connect(nil).Create("testdb", "message", "types", "Message", fileName); err != nil {
		t.Fatalf("Create() ошибка: %v", err)
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, want := range []string{"package types", "type Message struct", "ID", "uint64", "db:\"id\""} {
		if !strings.Contains(content, want) {
			t.Fatalf("сгенерированный файл не содержит %q:\n%s", want, content)
		}
	}

	// Ошибка получения информации о таблице.
	d = newImpl(&mockDialect{err: errors.New("сбой соединения")})
	if err := d.Connect(nil).Create("testdb", "message", "types", "Message", fileName); err == nil ||
		!strings.Contains(err.Error(), "ошибка получения информации о таблице") {
		t.Fatalf("ожидается ошибка получения информации, получено: %v", err)
	}

	// Ошибка генерации: неизвестный тип базы данных без переопределения.
	bad := tableInfo()
	bad.Columns[0].TypeGolang = d2sTypes.NewBase(d2sTypes.TUnknown)
	d = newImpl(&mockDialect{inf: bad})
	if err := d.Connect(nil).Create("testdb", "message", "types", "Message", fileName); err == nil ||
		!strings.Contains(err.Error(), "ошибка генерации структуры") {
		t.Fatalf("ожидается ошибка генерации, получено: %v", err)
	}

	// Ошибка открытия файла.
	d = newImpl(&mockDialect{inf: tableInfo()})
	if err := d.Connect(nil).Create("testdb", "message", "types", "Message",
		filepath.Join(t.TempDir(), "нет-такой-директории", "x.go")); err == nil ||
		!strings.Contains(err.Error(), "ошибка открытия файла") {
		t.Fatalf("ожидается ошибка открытия файла, получено: %v", err)
	}
}

// TestCreateWriteError Ошибка записи в файл: нулевой лимит размера файла RLIMIT_FSIZE.
func TestCreateWriteError(t *testing.T) {
	newImpl := func(m *mockDialect) *impl {
		d := New().(*impl)
		d.dialect = m

		return d
	}
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_FSIZE, &lim); err != nil {
		t.Skipf("RLIMIT_FSIZE не поддерживается: %v", err)
	}
	// Лимит процесса: восстанавливается по завершении теста.
	if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &syscall.Rlimit{Cur: 0, Max: lim.Max}); err != nil {
		t.Skipf("не удалось установить RLIMIT_FSIZE: %v", err)
	}
	defer func() {
		if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &lim); err != nil {
			t.Errorf("не удалось восстановить RLIMIT_FSIZE: %v", err)
		}
	}()

	d := newImpl(&mockDialect{inf: &d2sTypes.TableInfo{
		Database: "testdb", Table: "message",
		Columns: []*d2sTypes.ColumnInfo{
			{Name: "id", TypeGolang: d2sTypes.NewBase(d2sTypes.TUint64), TypeFull: "BIGINT(20)"},
		},
	}})
	err := d.Connect(nil).Create("testdb", "message", "types", "Message", filepath.Join(t.TempDir(), "x.go"))
	t.Logf("DEBUG err = %v", err)
	if err == nil || !strings.Contains(err.Error(), "ошибка записи в файл") {
		t.Fatalf("ожидается ошибка записи в файл, получено: %v", err)
	}
}
