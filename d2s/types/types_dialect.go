package types // import "gopkg.in/webnice/d2s.v1/d2s/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"database/sql"
)

// Dialect abstracts the details of specific SQL dialects
type Dialect interface {
	// DatabaseCurrent Получение имени текущей базы данных выбранной через DSN
	DatabaseCurrent(db *sql.DB) (ret string, err error)

	// TableInfo Получение информации о структуре таблицы
	TableInfo(db *sql.DB, databaseName string, tableName string) (ret *TableInfo, err error)
}
