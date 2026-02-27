package main

import "gopkg.in/alecthomas/kingpin.v2"

// Args Аргументы командной строки.
type Args struct {
	Debug     bool               // Debug flag.
	Driver    string             // Driver of database.
	Dsn       string             // Database source name (DSN).
	Database  string             // Database name.
	Table     string             // Table name.
	Package   string             // Package name.
	Structure string             // Structure name.
	File      string             // Имя файла для создания структуры golang.
	Create    *kingpin.CmdClause // Создает новый файл миграции со следующей версией.
}
