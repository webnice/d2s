package d2s // import "gopkg.in/webnice/d2s.v1/d2s"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"database/sql"

	d2sTypes "gopkg.in/webnice/d2s.v1/d2s/types"
)

// Interface is an interface of package
type Interface interface {
	// Debug sets a debug mode
	Debug(d bool) Interface

	// Dialect sets the SQL dialect
	Dialect(sqlDialect string) error

	// Connect Соединение с базой данных для получения информации
	Connect(db *sql.DB) Interface

	// Create structure from table
	Create(databaseName string, tableName string, packageName string, structureName string, fileName string) error
}

// impl is an implementation of package
type impl struct {
	debug   bool             // =true - Режим отладки
	dialect d2sTypes.Dialect // Диалект базы данных
	db      *sql.DB          // Объект соединения с базой данных
}
