package log

import "github.com/webnice/dic"

// Коды ошибок.
const (
	eLogPanicException           uint8 = iota + 1 // 001
	eHandlerAlreadySubscribed                     // 002
	eHandlerSubscriptionNotFound                  // 003
)

// Текстовые значения кодов ошибок на основном языке приложения.
const (
	cLogPanicException           = "Вызов обработчика сообщений лога %q прервана паникой:\n%v\n%s."
	cHandlerAlreadySubscribed    = "Обработчик логов %q уже зарегистрирован."
	cHandlerSubscriptionNotFound = "Регистрация обработчика логов %q не найдена."
)

type Error struct {
	dic.Errors

	// LogPanicException Вызов обработчика сообщений лога ... прервана паникой: ... ....
	LogPanicException dic.IError

	// HandlerAlreadySubscribed Обработчик логов ... уже зарегистрирован.
	HandlerAlreadySubscribed dic.IError

	// HandlerSubscriptionNotFound Регистрация обработчика логов ... не найдена.
	HandlerSubscriptionNotFound dic.IError
}

var (
	errSingleton = &Error{
		Errors:                      dic.Error(),
		LogPanicException:           dic.NewError(cLogPanicException, "обработчик логов", "паника", "стек вызовов").CodeU8().Set(eLogPanicException),
		HandlerAlreadySubscribed:    dic.NewError(cHandlerAlreadySubscribed, "обработчик логов").CodeU8().Set(eHandlerAlreadySubscribed),
		HandlerSubscriptionNotFound: dic.NewError(cHandlerSubscriptionNotFound, "обработчик логов").CodeU8().Set(eHandlerSubscriptionNotFound),
	}

	// Errors Справочник ошибок.
	Errors = func() *Error { return errSingleton }
)
