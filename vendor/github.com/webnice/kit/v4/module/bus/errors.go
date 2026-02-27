package bus

import "github.com/webnice/dic"

// Коды ошибок.
const (
	eDatabusRecursivePointer      uint8 = iota + 1 // 001
	eDatabusPanicException                         // 002
	eDatabusSubscribeNotFound                      // 003
	eDatabusInternalError                          // 004
	eDatabusNotSubscribersForType                  // 005
	eDatabusObjectIsNil                            // 006
	eDatabusPoolInternalError                      // 007
)

const (
	cDatabusRecursivePointer      = "Не возможно определить тип рекурсивного указателя: %q."
	cDatabusPanicException        = "Работа с подпиской потребителя, в шине данных, прервана паникой:\n%v\n%s."
	cDatabusSubscribeNotFound     = "Потребитель данных %q не был подписан на шину данных."
	cDatabusInternalError         = "Внутренняя ошибка шины данных: %s."
	cDatabusNotSubscribersForType = "Отсутствуют потребители данных для типа данных %q."
	cDatabusObjectIsNil           = "Передан nil объект."
	cDatabusPoolInternalError     = "Бассейн объектов вернул не корректный объект."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// DatabusRecursivePointer Не возможно определить тип рекурсивного указателя: ...
	DatabusRecursivePointer dic.IError

	// DatabusPanicException Работа с подпиской потребителя, в шине данных, прервана паникой:\n...\n...
	DatabusPanicException dic.IError

	// DatabusSubscribeNotFound Потребитель данных ... не был подписан на шину данных.
	DatabusSubscribeNotFound dic.IError

	// DatabusInternalError Внутренняя ошибка шины данных: ...
	DatabusInternalError dic.IError

	// DatabusNotSubscribersForType Отсутствуют потребители данных для типа данных ...
	DatabusNotSubscribersForType dic.IError

	// DatabusObjectIsNil Передан nil объект.
	DatabusObjectIsNil dic.IError

	// DatabusPoolInternalError Бассейн объектов вернул не корректный объект.
	DatabusPoolInternalError dic.IError
}

var errSingleton = &Error{
	Errors:                       dic.Error(),
	DatabusRecursivePointer:      dic.NewError(cDatabusRecursivePointer, "указатель").CodeU8().Set(eDatabusRecursivePointer),
	DatabusPanicException:        dic.NewError(cDatabusPanicException, "паника", "стек вызовов").CodeU8().Set(eDatabusPanicException),
	DatabusSubscribeNotFound:     dic.NewError(cDatabusSubscribeNotFound, "потребитель").CodeU8().Set(eDatabusSubscribeNotFound),
	DatabusInternalError:         dic.NewError(cDatabusInternalError, "ошибка").CodeU8().Set(eDatabusInternalError),
	DatabusNotSubscribersForType: dic.NewError(cDatabusNotSubscribersForType, "тип").CodeU8().Set(eDatabusNotSubscribersForType),
	DatabusObjectIsNil:           dic.NewError(cDatabusObjectIsNil).CodeU8().Set(eDatabusObjectIsNil),
	DatabusPoolInternalError:     dic.NewError(cDatabusPoolInternalError).CodeU8().Set(eDatabusPoolInternalError),
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
