package cli

import "github.com/webnice/dic"

const (
	cUnexpectedError    = "Неожиданная ошибка."
	cHelpDisplayed      = "Отображение помощи по командам, аргументам и флагам командной строки."
	cRequiredCommand    = "Не указана обязательная команда командной строки."
	cRequiredFlag       = "Не указан обязательный флаг командной строки."
	cUnknownCommand     = "Неизвестная команда командной строки."
	cUnknownArgument    = "Неизвестный аргумент командной строки."
	cNotCorrectArgument = "Не верное значение или тип аргумента, флага или параметра."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// UnexpectedError Неожиданная ошибка.
	UnexpectedError dic.IError

	// HelpDisplayed Отображение помощи по командам, аргументам и флагам командной строки.
	HelpDisplayed dic.IError

	// RequiredCommand Не указана обязательная команда командной строки.
	RequiredCommand dic.IError

	// RequiredFlag Не указан обязательный флаг командной строки.
	RequiredFlag dic.IError

	// UnknownCommand Неизвестная команда командной строки.
	UnknownCommand dic.IError

	// UnknownArgument Неизвестный аргумент командной строки.
	UnknownArgument dic.IError

	// NotCorrectArgument Не верное значение или тип аргумента, флага или параметра.
	NotCorrectArgument dic.IError
}

var (
	errSingleton = &Error{
		Errors:             dic.Error(),
		UnexpectedError:    dic.NewError(cUnexpectedError),
		HelpDisplayed:      dic.NewError(cHelpDisplayed),
		RequiredCommand:    dic.NewError(cRequiredCommand),
		RequiredFlag:       dic.NewError(cRequiredFlag),
		UnknownCommand:     dic.NewError(cUnknownCommand),
		UnknownArgument:    dic.NewError(cUnknownArgument),
		NotCorrectArgument: dic.NewError(cNotCorrectArgument),
	}

	// Errors Справочник ошибок.
	Errors = func() *Error { return errSingleton }
)
