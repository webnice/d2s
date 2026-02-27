package types

// TableInfo Информация о таблице.
type TableInfo struct {
	Database string        // Название базы данных.
	Table    string        // Название таблицы.
	Package  string        // Название пакета.
	Struct   string        // Название структуры.
	Comment  string        // Комментарий к структуре, он же комментарий к таблице.
	Columns  []*ColumnInfo // Колонки таблицы.
}

// ColumnInfo Информация о колонке.
type ColumnInfo struct {
	Name            ColumnName // [COLUMN_NAME] ............... Название колонки.
	Default         *string    // [COLUMN_DEFAULT] ............ Значение по умолчанию для колонки.
	TypeSimple      string     // [DATA_TYPE] ................. Тип колонки, простой вариант без расширения.
	TypeFull        string     // [COLUMN_TYPE] ............... Полный тип колонки используемый в запросе создания таблицы.
	Comment         string     // [COLUMN_COMMENT] ............ Комментарий к таблице.
	Size            *uint64    // [CHARACTER_MAXIMUM_LENGTH] .. Максимальная длинна данных колонки.
	Precision       *uint64    // [NUMERIC_PRECISION] ......... Точность.
	Scale           *uint64    // [NUMERIC_SCALE] ............. Масштаб (количество знаков после запятой).
	IsNullable      bool       // [IS_NULLABLE] ............... true - Значение может быть NULL.
	IsUnsigned      bool       // [COLUMN_TYPE] ............... true - Колонка UNSIGNED.
	IsPrimary       bool       // [COLUMN_KEY] ................ Колонка является первичным ключём.
	IsAutoIncrement bool       // [EXTRA] ..................... Если колонка имеет атрибут auto_increment.
	TypeGolang      *Base      // ............................. Тип данных golang
}
