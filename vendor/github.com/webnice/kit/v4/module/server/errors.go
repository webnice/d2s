package server

import "github.com/webnice/dic"

// Все ошибки определены как константы. Коды ошибок приложения:

// Коды ошибок.
const (
	eServerByIdNotFound uint8 = iota + 1
	eServersAddedNotEqualLaunched
	eModulePanicException
	eBeforeExitWithError
	eServerByIdAlreadyStarted
	eServerByIdNotStarted
)

// Текстовые значения кодов ошибок на основном языке приложения.
const (
	cServerByIdNotFound           = "Конфигурация сервера с идентификатором %q не найдена."
	cServersAddedNotEqualLaunched = "Запуск серверов завершился ошибкой, добавлено серверов: %d, запущено серверов: %d."
	cModulePanicException         = "Выполнение модуля %q прервано паникой:\n%v\n%s."
	cBeforeExitWithError          = "Функция Before() завершилась с ошибкой: %s"
	cServerByIdAlreadyStarted     = "Сервер с идентификатором %q уже запущен."
	cServerByIdNotStarted         = "Сервер с идентификатором %q не был запущен."
)

type Error struct {
	dic.Errors

	// ServerByIdNotFound Конфигурация сервера с идентификатором ... не найдена.
	ServerByIdNotFound dic.IError

	// ServersAddedNotEqualLaunched Добавлено ... серверов, запущено ... серверов.
	ServersAddedNotEqualLaunched dic.IError

	// ModulePanicException Выполнение модуля ... прервано паникой: ...
	ModulePanicException dic.IError

	// BeforeExitWithError Функция Before() завершилась с ошибкой: ...
	BeforeExitWithError dic.IError

	// ServerByIdAlreadyStarted Сервер с идентификатором ... уже запущен.
	ServerByIdAlreadyStarted dic.IError

	// ServerByIdNotStarted Сервер с идентификатором ... не был запущен.
	ServerByIdNotStarted dic.IError
}

var (
	errSingleton = &Error{
		Errors:                       dic.Error(),
		ServerByIdNotFound:           dic.NewError(cServerByIdNotFound, "идентификатор").CodeU8().Set(eServerByIdNotFound),
		ServersAddedNotEqualLaunched: dic.NewError(cServersAddedNotEqualLaunched, "добавлено", "запущено").CodeU8().Set(eServersAddedNotEqualLaunched),
		ModulePanicException:         dic.NewError(cModulePanicException, "название модуля", "стек паники").CodeU8().Set(eModulePanicException),
		BeforeExitWithError:          dic.NewError(cBeforeExitWithError, "ошибка").CodeU8().Set(eBeforeExitWithError),
		ServerByIdAlreadyStarted:     dic.NewError(cServerByIdAlreadyStarted, "идентификатор").CodeU8().Set(eServerByIdAlreadyStarted),
		ServerByIdNotStarted:         dic.NewError(cServerByIdNotStarted, "идентификатор").CodeU8().Set(eServerByIdNotStarted),
	}

	// Errors Справочник ошибок.
	Errors = func() *Error { return errSingleton }
)
