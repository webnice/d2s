package d2s // import "gopkg.in/webnice/d2s.v1/d2s"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// Interface is an interface of package
type Interface interface {
	// Debug sets a debug mode
	Debug(d bool) Interface

	// Dialect sets the SQL dialect
	Dialect(sqlDialect string) error

	// Run runs a command
	Run(command string, tableName string, structureName string, fileName string) error
}

// impl is an implementation of package
type impl struct {
	debug   bool
	dialect SQLDialect
}

// Sqlite3Dialect struct
type Sqlite3Dialect struct{}

// MySQLDialect struct
type MySQLDialect struct{}

// ClickHouseDialect struct
type ClickHouseDialect struct{}

// RedshiftDialect struct
type RedshiftDialect struct{}

// TiDBDialect struct
type TiDBDialect struct{}

// PostgresDialect struct
type PostgresDialect struct{}

// SQLDialect abstracts the details of specific SQL dialects
type SQLDialect interface {
}
