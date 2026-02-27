package cpy

import "github.com/webnice/dic"

const (
	cCopyToObjectUnaddressable = "Объект назначения не адресуемый."
	cCopyFromObjectInvalid     = "Объект источника не адресуемый."
	cTypeMapNotEqual           = "Типы данных карты не эквивалентны."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// CopyToObjectUnaddressable Объект назначения не адресуемый.
	CopyToObjectUnaddressable dic.IError

	// CopyFromObjectInvalid Объект источника не адресуемый.
	CopyFromObjectInvalid dic.IError

	// TypeMapNotEqual Типы данных карты не эквивалентны.
	TypeMapNotEqual dic.IError
}

var errSingleton = &Error{
	Errors:                    dic.Error(),
	CopyToObjectUnaddressable: dic.NewError(cCopyToObjectUnaddressable),
	CopyFromObjectInvalid:     dic.NewError(cCopyFromObjectInvalid),
	TypeMapNotEqual:           dic.NewError(cTypeMapNotEqual),
}
