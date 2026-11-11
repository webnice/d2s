package d2s

import (
	"database/sql"

	d2sTypes "github.com/webnice/d2s/v2/types"
)

// Interface Интерфейс пакета.
type Interface interface {
	// Debug Установка режима отладки.
	Debug(d bool) Interface

	// Dialect Установка диалекта SQL.
	Dialect(sqlDialect string) error

	// Types Переопределение типов полей структуры.
	// Список "поле1:тип;поле2:тип", где поле — имя колонки базы данных,
	// тип — стандартный тип Golang или тип nix/v2.
	Types(overrideList string) error

	// Connect Соединение с базой данных для получения информации.
	Connect(db *sql.DB) Interface

	// Create Создание структуры по таблице базы данных.
	Create(databaseName string, tableName string, packageName string, structureName string, fileName string) error
}

// impl Реализация пакета.
type impl struct {
	debug    bool              // =true - Режим отладки.
	dialect  d2sTypes.Dialect  // Диалект базы данных.
	db       *sql.DB           // Объект соединения с базой данных.
	override map[string]string // Переопределение типов полей структуры.
}
