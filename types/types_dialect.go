// Package types
package types

import "database/sql"

// Dialect abstracts the details of specific SQL dialects
type Dialect interface {
	// DatabaseCurrent Получение имени текущей базы данных выбранной через DSN
	DatabaseCurrent(db *sql.DB) (ret string, err error)

	// TableInfo Получение информации о структуре таблицы
	TableInfo(db *sql.DB, databaseName string, tableName string) (ret *TableInfo, err error)
}
