package main

import "github.com/alecthomas/kingpin/v2"

func args() (cmd string, args *Args) {
	args = new(Args)
	kingpin.CommandLine.Help = `Утилита создаёт golang структуру по информации о таблице базы данных`

	// Глобальные флаги.
	kingpin.Flag(`debug`, `Включение режима отладки.
Переопределяет значение по умолчанию флага из переменной окружения с именем 'DB2STRUCT_DEBUG'`).
		Envar("DB2STRUCT_DEBUG").
		Default("false").
		Short('d').
		BoolVar(&args.Debug)
	kingpin.Flag(`drv`, `Драйвер базы данных.
Переопределяет значение по умолчанию флага из переменной окружения с именем 'DB2STRUCT_DRV'`).
		Envar("DB2STRUCT_DRV").
		Default("mysql").
		Short('b').
		StringVar(&args.Driver)
	kingpin.Flag(`dsn`, `Имя источника данных (DSN).
Переопределяет значение по умолчанию флага из переменной окружения с именем 'DB2STRUCT_DSN'`).
		Envar(`DB2STRUCT_DSN`).
		Default(`root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true`).
		Short('u').
		StringVar(&args.Dsn)

	// Команды с аргументами.
	args.Create = kingpin.Command(`create`, `Создание структуры с указанным именем по данным
из базы данных и запись результата в указанный файл.`)
	args.Create.Arg(`database`, `Имя базы данных. Если не указано, определяется из строки DSN или соединения.`).
		StringVar(&args.Database)
	args.Create.Arg(`table`, `Имя таблицы.`).
		StringVar(&args.Table)
	args.Create.Arg(`package`, `Имя пакета.`).
		StringVar(&args.Package)
	args.Create.Arg(`struct`, `Имя создаваемой структуры.`).
		StringVar(&args.Structure)
	args.Create.Arg(`file`, `Имя создаваемого файла со структурой таблицы.`).
		StringVar(&args.File)
	args.Create.Arg(`types`, `Переопределение типов данных полей структуры.
Формат: "поле1:тип;поле2:тип", где поле — имя колонки базы данных,
тип — стандартный тип данных golang или тип данных библиотеки github.com/webnice/nix/v2.
Переопределяет значение по умолчанию, только если поле базы данных указано и совпадает.`).
		StringVar(&args.Types)
	cmd = kingpin.Parse()

	return
}

func argUsage() {
	kingpin.Usage()
}
