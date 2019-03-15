package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func args() (cmd string, args *Args) {
	args = new(Args)
	kingpin.CommandLine.Help = `The utility creates a golang structure from information about a database table`

	// Global flags
	kingpin.Flag(`debug`, `Sets debug mode.
Overrides the default value for a flag from an environment variable by name 'DB2STRUCT_DEBUG'`).
		Envar("DB2STRUCT_DEBUG").
		Default("false").
		Short('d').
		BoolVar(&args.Debug)
	kingpin.Flag(`drv`, `Driver of database.
Overrides the default value for a flag from an environment variable by name 'DB2STRUCT_DRV'`).
		Envar("DB2STRUCT_DRV").
		Default("mysql").
		Short('b').
		StringVar(&args.Driver)
	kingpin.Flag(`dsn`, `Database source name (DSN).
Overrides the default value for a flag from an environment variable by name 'DB2STRUCT_DSN'`).
		Envar(`DB2STRUCT_DSN`).
		Default(`root@unix(/var/run/mysql/mysql.sock)/test?parseTime=true`).
		Short('u').
		StringVar(&args.Dsn)

	// Commands with args
	args.Create = kingpin.Command(`create`, `Create a structure file`)
	args.Create.Arg(`table`, `Database name and table name through point.`).
		StringVar(&args.Table)
	args.Create.Arg(`struct`, `Name of created structure.`).
		StringVar(&args.Structure)
	args.Create.Arg(`file`, `The name of the file being created with the table structure.`).
		StringVar(&args.File)
	cmd = kingpin.Parse()

	return
}

func argUsage() {
	kingpin.Usage()
}
