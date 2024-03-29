// Package types
package types

// Base type.
type Base struct {
	Simple     string // Простейший базовый тип.
	IsNullable bool   // Флаг nullable, =true значение может быть NULL.
}

// NewBase Create new object type Base.
func NewBase(simple string) *Base {
	var bt = &Base{simple, false}
	return bt
}

// Nullable Set nullable flag.
func (bt *Base) Nullable(isNullable bool) { bt.IsNullable = isNullable }

// String Return type as string.
func (bt *Base) String() string {
	if _, ok := typesMap[bt.Simple]; !ok {
		return ""
	}
	return typesMap[bt.Simple].Nullable(bt.IsNullable)
}
