package main

import (
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"
)

// TestLog Ссылка на менеджер логирования.
func TestLog(t *testing.T) {
	if log() == nil {
		t.Fatal("log() вернул nil")
	}
}

// TestArgs Разбор аргументов командной строки.
func TestArgs(t *testing.T) {
	oldCommandLine := kingpin.CommandLine
	defer func() { kingpin.CommandLine = oldCommandLine }()
	kingpin.CommandLine = kingpin.New("db2struct", "")
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"db2struct",
		"create", "mydb", "message", "types", "Message", "message_model.go", "id:uint64;payload:nix.RawMessage"}
	cmd, args := args()
	if cmd != "create" {
		t.Fatalf("cmd = %q; ожидается create", cmd)
	}
	if args.Database != "mydb" || args.Table != "message" || args.Package != "types" ||
		args.Structure != "Message" || args.File != "message_model.go" ||
		args.Types != "id:uint64;payload:nix.RawMessage" {
		t.Fatalf("args = %#v", args)
	}
	// Справка по использованию.
	argUsage()
}

// TestMainMain Запуск main() во вспомогательном процессе того же тестового бинарника.
// Менеджер логирования kit не завершает процесс — main() отрабатывает до конца.
func TestMainMain(t *testing.T) {
	oldCommandLine := kingpin.CommandLine
	defer func() { kingpin.CommandLine = oldCommandLine }()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for _, tc := range []struct {
		name string
		args []string
		dsn  string
	}{
		{"mysql без соединения", []string{"create", "", "message", "types", "Message", "x.go", "id:uint64"}, ""},
		{"postgres не реализован", []string{"--drv=postgres", "create", "", "message", "types", "x.go", ""}, ""},
		{"sqlite3 не реализован", []string{"--drv=sqlite3", "create", "", "message", "types", "x.go", ""}, ""},
		{"clickhouse не реализован", []string{"--drv=clickhouse", "create", "", "message", "types", "x.go", ""}, ""},
		{"redshift не реализован", []string{"--drv=redshift", "create", "", "message", "types", "x.go", ""}, ""},
		{"tidb не реализован", []string{"--drv=tidb", "create", "", "message", "types", "x.go", ""}, ""},
		{"неизвестный драйвер", []string{"--drv=oracle", "create", "", "message", "types", "x.go", ""}, ""},
		{"ошибка переопределения типов", []string{"create", "", "message", "types", "Message", "x.go", "id:uint128"}, ""},
		{"некорректный DSN", []string{"create", "", "message", "types", "Message", "x.go", ""}, "://некорректный"},
		{"пустое имя таблицы", []string{"create", "", "", "types", "Message", "x.go", ""}, ""},
		{"команда help", []string{"help"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Свежий объект приложения kingpin: повторная регистрация команд запрещена.
			// Terminate отключён: kingpin не должен завершать тестовый процесс.
			app := kingpin.New("db2struct", "")
			app.Terminate(func(int) {})
			kingpin.CommandLine = app
			if tc.dsn != "" {
				t.Setenv("DB2STRUCT_DSN", tc.dsn)
			} else {
				t.Setenv("DB2STRUCT_DSN", "root:root@tcp(127.0.0.1:1)/test?parseTime=true")
			}
			os.Args = append([]string{"db2struct"}, tc.args...)
			// Менеджер логирования не завершает процесс: main() отрабатывает до конца.
			main()
		})
	}
}
