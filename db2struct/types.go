package main

import "github.com/alecthomas/kingpin/v2"

// Args Аргументы командной строки.
type Args struct {
	Debug     bool               // Флаг режима отладки.
	Driver    string             // Драйвер базы данных.
	Dsn       string             // Имя источника данных (DSN).
	Database  string             // Имя базы данных.
	Table     string             // Имя таблицы.
	Package   string             // Имя пакета.
	Structure string             // Имя структуры.
	File      string             // Имя файла для создания структуры Golang.
	Types     string             // Переопределение типов данных полей структуры.
	Create    *kingpin.CmdClause // Команда создания структуры по таблице базы данных.
}
