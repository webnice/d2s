package server

import "github.com/webnice/kit/v4/module/ambry"

const (
	errInternalServerError uint32 = iota + 1
	errMethodNotAllowed
)

// InterfaceWebError Интерфейс ошибок ВЕБ сервера.
type InterfaceWebError interface {
	// Reset Сброс всех установленных ошибок.
	Reset()

	// InternalServerError Установка и получение значения ошибки "Внутренняя ошибка сервера".
	InternalServerError(err error) error

	// MethodNotAllowed Установка и получение значения ошибки "Метод не разрешён".
	MethodNotAllowed(err error) error
}

// Объект сущности, реализующий интерфейс InterfaceWebError.
type implWebError struct {
	parent *impl           // Адрес объекта родительской сущности (parent), интерфейс Interface.
	errors ambry.Interface // Интерфейс хранения ошибок.
}
