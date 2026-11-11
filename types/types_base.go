package types

// Base Базовый тип.
type Base struct {
	Simple     string // Простейший базовый тип.
	IsNullable bool   // Флаг nullable, =true значение может быть NULL.
}

// NewBase Конструктор объекта базового типа.
func NewBase(simple string) *Base {
	var bt = &Base{simple, false}
	return bt
}

// Nullable Установка флага nullable.
func (bt *Base) Nullable(isNullable bool) { bt.IsNullable = isNullable }

// String Строковое представление типа.
func (bt *Base) String() string {
	if _, ok := typesMap[bt.Simple]; !ok {
		return ""
	}
	return typesMap[bt.Simple].Nullable(bt.IsNullable)
}
