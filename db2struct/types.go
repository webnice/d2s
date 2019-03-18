package main

import "gopkg.in/alecthomas/kingpin.v2"

// Args The arguments
type Args struct {
	Debug     bool               // Debug flag
	Driver    string             // Driver of database
	Dsn       string             // Database source name (DSN)
	Database  string             // Database name
	Table     string             // Table name
	Package   string             // Package name
	Structure string             // Structure name
	File      string             // Name of file for create golang structure
	Create    *kingpin.CmdClause // Creates new migration file with next version
}
