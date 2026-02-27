package cfg

import "github.com/webnice/dic"

// Коды особенных ошибок.
const (
	eApplicationPanicException uint8 = 255 // 255 - приложение завершилось из-за прерывания по исключению - panic.
	eApplicationUnknownError   uint8 = 254 // 254 - неожиданная ошибка приложения.
	eApplicationHelpDisplayed  uint8 = 253 // 253 - приложение завершилось корректно с отображением помощи по CLI.
)

// Коды ошибок.
const (
	eApplicationVersion                        uint8 = iota + 1 // 001
	eApplicationMainFuncNotFound                                // 002
	eApplicationMainFuncAlreadyRegistered                       // 003
	eApplicationRegistrationUnknownObject                       // 004
	eComponentIsNull                                            // 005
	eComponentRegistrationProhibited                            // 006
	eComponentRegistrationError                                 // 007
	eComponentPreferencesCallBeforeCompleting                   // 008
	eComponentPanicException                                    // 009
	eComponentRunlevelError                                     // 010
	eComponentRulesError                                        // 011
	eRunlevelCantLessCurrentLevel                               // 012
	eInitLogging                                                // 013
	eComponentConflict                                          // 014
	eComponentRequires                                          // 015
	eComponentInitiateTimeout                                   // 016
	eComponentInitiateExecution                                 // 017
	eComponentInitiatePanicException                            // 018
	eComponentDoExecution                                       // 019
	eComponentDoPanicException                                  // 020
	eComponentDoUnknownError                                    // 021
	eComponentFinalizeExecution                                 // 022
	eComponentFinalizePanicException                            // 023
	eComponentFinalizeUnknownError                              // 024
	eComponentFinalizeWarning                                   // 025
	eRunlevelSubscribeUnsubscribeNilFunction                    // 026
	eRunlevelAlreadySubscribedFunction                          // 027
	eRunlevelSubscriptionNotFound                               // 028
	eRunlevelSubscriptionPanicException                         // 029
	eCommandLineArgumentRequired                                // 030
	eCommandLineArgumentUnknown                                 // 031
	eCommandLineArgumentNotCorrect                              // 032
	eCommandLineRequiredFlag                                    // 033
	eCommandLineUnexpectedError                                 // 034
	eConfigurationBootstrap                                     // 035
	eGetCurrentUser                                             // 036
	eCantChangeWorkDirectory                                    // 037
	ePidExistsAnotherProcessOfApplication                       // 038
	ePidFileError                                               // 039
	eDatabusRecursivePointer                                    // 040
	eDatabusPanicException                                      // 041
	eDatabusSubscribeNotFound                                   // 042
	eDatabusInternalError                                       // 043
	eDatabusNotSubscribersForType                               // 044
	eDatabusObjectIsNil                                         // 045
	eConfigurationApplicationProhibited                         // 046
	eConfigurationApplicationObject                             // 047
	eConfigurationApplicationPanic                              // 048
	eConfigurationFileNotFound                                  // 049
	eConfigurationPermissionDenied                              // 050
	eConfigurationUnexpectedMistakeFileAccess                   // 051
	eConfigurationFileIsDirectory                               // 052
	eConfigurationFileReadingError                              // 053
	eConfigurationSetDefault                                    // 054
	eConfigurationSetDefaultValue                               // 055
	eConfigurationSetDefaultPanic                               // 056
	eConfigurationObjectNotFound                                // 057
	eConfigurationObjectIsNotStructure                          // 058
	eConfigurationObjectIsNil                                   // 059
	eConfigurationObjectIsNotValid                              // 060
	eConfigurationObjectIsNotAddress                            // 061
	eConfigurationObjectCopy                                    // 062
	eConfigurationCallbackAlreadyRegistered                     // 063
	eConfigurationCallbackSubscriptionNotFound                  // 064
)

// Текстовые значения кодов ошибок.
const (
	cApplicationPanicException                 = "Выполнение приложения прервано паникой:\n%v\n%s."
	cApplicationUnknownError                   = "%s"
	cApplicationHelpDisplayed                  = "%s"
	cApplicationVersion                        = "Версии приложения содержит ошибку: %s."
	cApplicationMainFuncNotFound               = "Не определена основная функция приложения."
	cApplicationMainFuncAlreadyRegistered      = "Основная функция приложения уже зарегистрирована."
	cApplicationRegistrationUnknownObject      = "Регистрация не известного компонента, объекта или модуля: %q."
	cComponentIsNull                           = "В качестве объекта компоненты передан nil."
	cComponentRegistrationProhibited           = "Регистрация компонентов запрещена. Компонента %q не зарегистрирована."
	cComponentRegistrationError                = "Регистрация компоненты %q завершилась ошибкой: %s."
	cComponentPreferencesCallBeforeCompleting  = "Опрос настроек компонентов вызван до завершения регистрации компонентов."
	cComponentPanicException                   = "Выполнение компоненты %q прервано паникой:\n%v\n%s."
	cComponentRunlevelError                    = "Уровень запуска (runlevel), для компоненты %q указан %d, необходимо указать уровень равный 0, либо в интервале от 10 до 65534 включительно."
	cComponentRulesError                       = "Правила %q для компоненты %q содержат ошибку: %s."
	cRunlevelCantLessCurrentLevel              = "Новый уровень работы приложения (%d) не может быть меньше текущего уровня работы приложения (%d)."
	cInitLogging                               = "Критическая ошибка в модуле менеджера логирования: %s."
	cComponentConflict                         = "Компонента %q конфликтует с компонентой %q."
	cComponentRequires                         = "Компонента %q имеет не удовлетворённую зависимость %q."
	cComponentInitiateTimeout                  = "Превышено время ожидание выполнения функции Initiate() компоненты %q."
	cComponentInitiateExecution                = "Выполнение функции Initiate() компоненты %q завершено с ошибкой: %s."
	cComponentInitiatePanicException           = "Выполнение функции Initiate() компоненты %q прервано паникой:\n%v\n%s."
	cComponentDoExecution                      = "Выполнение функции Do() компоненты %q завершено с ошибкой: %s."
	cComponentDoPanicException                 = "Выполнение функции Do() компоненты %q прервано паникой:\n%v\n%s."
	cComponentDoUnknownError                   = "Выполнение функций Do() завершилось ошибкой: %s."
	cComponentFinalizeExecution                = "Выполнение функции Finalize() компоненты %q завершено с ошибкой: %s."
	cComponentFinalizePanicException           = "Выполнение функции Finalize() компоненты %q прервано паникой:\n%v\n%s."
	cComponentFinalizeUnknownError             = "Выполнение функций Finalize() завершилось ошибкой: %s."
	cComponentFinalizeWarning                  = "Выполнение функций Finalize() компоненты %q, длится дольше отведённого времени (%s)."
	cRunlevelSubscribeUnsubscribeNilFunction   = "Передана nil функция, подписка или отписка nil функции не возможна."
	cRunlevelAlreadySubscribedFunction         = "Функция %q уже подписана на получение событий изменения уровня работы приложения."
	cRunlevelSubscriptionNotFound              = "Не найдена подписка функции %q на события изменения уровня работы приложения."
	cRunlevelSubscriptionPanicException        = "Вызов функции подписчика на событие изменения уровня работы приложения, прервано паникой:\n%v\n%s."
	cCommandLineArgumentRequired               = "Требуется указать обязательную команду, аргумент или флаг командной строки: %s."
	cCommandLineArgumentUnknown                = "Неизвестная команда, аргумент или флаг командной строки: %s."
	cCommandLineArgumentNotCorrect             = "Не верное значение или тип аргумента, флага или параметра: %s."
	cCommandLineRequiredFlag                   = "Не указан один или несколько обязательных флагов: %s."
	cCommandLineUnexpectedError                = "Не предвиденная ошибка библиотеки командного интерфейса приложения: %s; %s."
	cConfigurationBootstrap                    = "Ошибка начально bootstrap конфигурации приложения: %s."
	cGetCurrentUser                            = "Не удалось загрузить данные о текущем пользователе операционной системы: %s."
	cCantChangeWorkDirectory                   = "Не удалось сменить рабочую директорию приложения: %s."
	cPidExistsAnotherProcessOfApplication      = "Существует один или несколько работающих процессов приложения, измените PID файл или остановите экземпляры приложения, PID: %s."
	cPidFileError                              = "Ошибка работы с PID файлом %q: %s."
	cDatabusRecursivePointer                   = "Не возможно определить тип рекурсивного указателя: %q."
	cDatabusPanicException                     = "Работа с подпиской потребителя, в шине данных, прервана паникой:\n%v\n%s."
	cDatabusSubscribeNotFound                  = "Потребитель данных %q не был подписан на шину данных."
	cDatabusInternalError                      = "Внутренняя ошибка шины данных: %s."
	cDatabusNotSubscribersForType              = "Отсутствуют потребители данных для типа данных: %q."
	cDatabusObjectIsNil                        = "Передан nil объект."
	cConfigurationApplicationProhibited        = "Регистрация объектов конфигурации на текущем уровне работы приложения запрещена. Конфигурация %q не зарегистрирована."
	cConfigurationApplicationObject            = "Объект конфигурации приложения содержит ошибку: %s."
	cConfigurationApplicationPanic             = "Непредвиденная ошибка при регистрации объекта конфигурации.\nПаника: %v.\n%s"
	cConfigurationFileNotFound                 = "Указанного файла конфигурации %q не существует: %s."
	cConfigurationPermissionDenied             = "Отсутствует доступ к файлу конфигурации %q, ошибка: %s."
	cConfigurationUnexpectedMistakeFileAccess  = "Неожиданная ошибка доступа к файлу конфигурации %q: %s."
	cConfigurationFileIsDirectory              = "В качестве файла конфигурации указана директория: %s."
	cConfigurationFileReadingError             = "Чтение фала конфигурации %q прервано ошибкой: %s."
	cConfigurationSetDefault                   = "Установка значений по умолчанию, для переменных конфигурации, прервана ошибкой: %s."
	cConfigurationSetDefaultValue              = "Установка значения по умолчанию %q, для переменной конфигурации %q, прервана ошибкой: %s."
	cConfigurationSetDefaultPanic              = "Непредвиденная ошибка, при установке значений по умолчанию, объекта конфигурации.\nПаника: %v.\n%s"
	cConfigurationObjectNotFound               = "Объект конфигурации с типом %q не найден."
	cConfigurationObjectIsNotStructure         = "Переданный объект %q не является структурой."
	cConfigurationObjectIsNil                  = "Переданный объект, является nil объектом."
	cConfigurationObjectIsNotValid             = "Объект конфигурации с типом %q не инициализирован."
	cConfigurationObjectIsNotAddress           = "Объект конфигурации с типом %q передан не корректно. Необходимо передать адрес объекта."
	cConfigurationObjectCopy                   = "Копирование объекта конфигурации с типом %q прервано ошибкой: %s."
	cConfigurationCallbackAlreadyRegistered    = "Подписка функции обратного вызова на изменение конфигурации с типом %q для функции %q уже существует."
	cConfigurationCallbackSubscriptionNotFound = "Подписка функции обратного вызова на изменение конфигурации с типом %q для функции %q не существует."
)

var errSingleton = &Error{
	Errors:                                    dic.Error(),
	ApplicationPanicException:                 dic.NewError(cApplicationPanicException, "паника", "стек вызовов").CodeU8().Set(eApplicationPanicException),
	ApplicationUnknownError:                   dic.NewError(cApplicationUnknownError, "ошибка").CodeU8().Set(eApplicationUnknownError),
	ApplicationHelpDisplayed:                  dic.NewError(cApplicationHelpDisplayed, "помощь").CodeU8().Set(eApplicationHelpDisplayed),
	ApplicationVersion:                        dic.NewError(cApplicationVersion, "ошибка").CodeU8().Set(eApplicationVersion),
	ApplicationMainFuncNotFound:               dic.NewError(cApplicationMainFuncNotFound).CodeU8().Set(eApplicationMainFuncNotFound),
	ApplicationMainFuncAlreadyRegistered:      dic.NewError(cApplicationMainFuncAlreadyRegistered).CodeU8().Set(eApplicationMainFuncAlreadyRegistered),
	ApplicationRegistrationUnknownObject:      dic.NewError(cApplicationRegistrationUnknownObject, "объект").CodeU8().Set(eApplicationRegistrationUnknownObject),
	ComponentIsNull:                           dic.NewError(cComponentIsNull).CodeU8().Set(eComponentIsNull),
	ComponentRegistrationProhibited:           dic.NewError(cComponentRegistrationProhibited, "название компоненты").CodeU8().Set(eComponentRegistrationProhibited),
	ComponentRegistrationError:                dic.NewError(cComponentRegistrationError, "название компоненты", "ошибка").CodeU8().Set(eComponentRegistrationError),
	ComponentPreferencesCallBeforeCompleting:  dic.NewError(cComponentPreferencesCallBeforeCompleting).CodeU8().Set(eComponentPreferencesCallBeforeCompleting),
	ComponentPanicException:                   dic.NewError(cComponentPanicException, "название компоненты", "паника", "стек вызовов").CodeU8().Set(eComponentPanicException),
	ComponentRunlevelError:                    dic.NewError(cComponentRunlevelError, "название компоненты", "уровень").CodeU8().Set(eComponentRunlevelError),
	ComponentRulesError:                       dic.NewError(cComponentRulesError, "правила", "название компоненты", "ошибка").CodeU8().Set(eComponentRulesError),
	RunlevelCantLessCurrentLevel:              dic.NewError(cRunlevelCantLessCurrentLevel, "уровень", "уровень").CodeU8().Set(eRunlevelCantLessCurrentLevel),
	InitLogging:                               dic.NewError(cInitLogging, "ошибка").CodeU8().Set(eInitLogging),
	ComponentConflict:                         dic.NewError(cComponentConflict, "название компоненты", "название компоненты").CodeU8().Set(eComponentConflict),
	ComponentRequires:                         dic.NewError(cComponentRequires, "название компоненты", "зависимость").CodeU8().Set(eComponentRequires),
	ComponentInitiateTimeout:                  dic.NewError(cComponentInitiateTimeout, "название компоненты").CodeU8().Set(eComponentInitiateTimeout),
	ComponentInitiateExecution:                dic.NewError(cComponentInitiateExecution, "название компоненты", "ошибка").CodeU8().Set(eComponentInitiateExecution),
	ComponentInitiatePanicException:           dic.NewError(cComponentInitiatePanicException, "название компоненты", "паника", "стек вызовов").CodeU8().Set(eComponentInitiatePanicException),
	ComponentDoExecution:                      dic.NewError(cComponentDoExecution, "название компоненты", "ошибка").CodeU8().Set(eComponentDoExecution),
	ComponentDoPanicException:                 dic.NewError(cComponentDoPanicException, "название компоненты", "паника", "стек вызовов").CodeU8().Set(eComponentDoPanicException),
	ComponentDoUnknownError:                   dic.NewError(cComponentDoUnknownError, "ошибка").CodeU8().Set(eComponentDoUnknownError),
	ComponentFinalizeExecution:                dic.NewError(cComponentFinalizeExecution, "название компоненты", "ошибка").CodeU8().Set(eComponentFinalizeExecution),
	ComponentFinalizePanicException:           dic.NewError(cComponentFinalizePanicException, "название компоненты", "паника", "стек вызовов").CodeU8().Set(eComponentFinalizePanicException),
	ComponentFinalizeUnknownError:             dic.NewError(cComponentFinalizeUnknownError, "ошибка").CodeU8().Set(eComponentFinalizeUnknownError),
	ComponentFinalizeWarning:                  dic.NewError(cComponentFinalizeWarning, "название компоненты", "время").CodeU8().Set(eComponentFinalizeWarning),
	RunlevelSubscribeUnsubscribeNilFunction:   dic.NewError(cRunlevelSubscribeUnsubscribeNilFunction).CodeU8().Set(eRunlevelSubscribeUnsubscribeNilFunction),
	RunlevelAlreadySubscribedFunction:         dic.NewError(cRunlevelAlreadySubscribedFunction, "название функции").CodeU8().Set(eRunlevelAlreadySubscribedFunction),
	RunlevelSubscriptionNotFound:              dic.NewError(cRunlevelSubscriptionNotFound, "название функции").CodeU8().Set(eRunlevelSubscriptionNotFound),
	RunlevelSubscriptionPanicException:        dic.NewError(cRunlevelSubscriptionPanicException, "паника", "стек вызовов").CodeU8().Set(eRunlevelSubscriptionPanicException),
	CommandLineArgumentRequired:               dic.NewError(cCommandLineArgumentRequired, "параметр").CodeU8().Set(eCommandLineArgumentRequired),
	CommandLineArgumentUnknown:                dic.NewError(cCommandLineArgumentUnknown, "неизвестный параметр").CodeU8().Set(eCommandLineArgumentUnknown),
	CommandLineArgumentNotCorrect:             dic.NewError(cCommandLineArgumentNotCorrect, "параметр").CodeU8().Set(eCommandLineArgumentNotCorrect),
	CommandLineRequiredFlag:                   dic.NewError(cCommandLineRequiredFlag, "название флага").CodeU8().Set(eCommandLineRequiredFlag),
	CommandLineUnexpectedError:                dic.NewError(cCommandLineUnexpectedError, "ошибка", "ошибка").CodeU8().Set(eCommandLineUnexpectedError),
	ConfigurationBootstrap:                    dic.NewError(cConfigurationBootstrap, "ошибка").CodeU8().Set(eConfigurationBootstrap),
	GetCurrentUser:                            dic.NewError(cGetCurrentUser, "пользователь").CodeU8().Set(eGetCurrentUser),
	CantChangeWorkDirectory:                   dic.NewError(cCantChangeWorkDirectory, "директория").CodeU8().Set(eCantChangeWorkDirectory),
	PidExistsAnotherProcessOfApplication:      dic.NewError(cPidExistsAnotherProcessOfApplication, "идентификатор").CodeU8().Set(ePidExistsAnotherProcessOfApplication),
	PidFileError:                              dic.NewError(cPidFileError, "название файла", "ошибка").CodeU8().Set(ePidFileError),
	DatabusRecursivePointer:                   dic.NewError(cDatabusRecursivePointer, "указатель").CodeU8().Set(eDatabusRecursivePointer),
	DatabusPanicException:                     dic.NewError(cDatabusPanicException, "паника", "стек вызовов").CodeU8().Set(eDatabusPanicException),
	DatabusSubscribeNotFound:                  dic.NewError(cDatabusSubscribeNotFound, "потребитель").CodeU8().Set(eDatabusSubscribeNotFound),
	DatabusInternalError:                      dic.NewError(cDatabusInternalError, "ошибка").CodeU8().Set(eDatabusInternalError),
	DatabusNotSubscribersForType:              dic.NewError(cDatabusNotSubscribersForType, "тип данных").CodeU8().Set(eDatabusNotSubscribersForType),
	DatabusObjectIsNil:                        dic.NewError(cDatabusObjectIsNil).CodeU8().Set(eDatabusObjectIsNil),
	ConfigurationApplicationProhibited:        dic.NewError(cConfigurationApplicationProhibited, "конфигурация").CodeU8().Set(eConfigurationApplicationProhibited),
	ConfigurationApplicationObject:            dic.NewError(cConfigurationApplicationObject, "ошибка").CodeU8().Set(eConfigurationApplicationObject),
	ConfigurationApplicationPanic:             dic.NewError(cConfigurationApplicationPanic, "паника", "стек вызовов").CodeU8().Set(eConfigurationApplicationPanic),
	ConfigurationFileNotFound:                 dic.NewError(cConfigurationFileNotFound, "название файла", "ошибка").CodeU8().Set(eConfigurationFileNotFound),
	ConfigurationPermissionDenied:             dic.NewError(cConfigurationPermissionDenied, "название файла", "ошибка").CodeU8().Set(eConfigurationPermissionDenied),
	ConfigurationUnexpectedMistakeFileAccess:  dic.NewError(cConfigurationUnexpectedMistakeFileAccess, "название файла", "ошибка").CodeU8().Set(eConfigurationUnexpectedMistakeFileAccess),
	ConfigurationFileIsDirectory:              dic.NewError(cConfigurationFileIsDirectory, "название директории").CodeU8().Set(eConfigurationFileIsDirectory),
	ConfigurationFileReadingError:             dic.NewError(cConfigurationFileReadingError, "название файла", "ошибка").CodeU8().Set(eConfigurationFileReadingError),
	ConfigurationSetDefault:                   dic.NewError(cConfigurationSetDefault, "ошибка").CodeU8().Set(eConfigurationSetDefault),
	ConfigurationSetDefaultValue:              dic.NewError(cConfigurationSetDefaultValue, "значение", "переменная", "ошибка").CodeU8().Set(eConfigurationSetDefaultValue),
	ConfigurationSetDefaultPanic:              dic.NewError(cConfigurationSetDefaultPanic, "паника", "стек вызовов").CodeU8().Set(eConfigurationSetDefaultPanic),
	ConfigurationObjectNotFound:               dic.NewError(cConfigurationObjectNotFound, "тип").CodeU8().Set(eConfigurationObjectNotFound),
	ConfigurationObjectIsNotStructure:         dic.NewError(cConfigurationObjectIsNotStructure, "объект").CodeU8().Set(eConfigurationObjectIsNotStructure),
	ConfigurationObjectIsNil:                  dic.NewError(cConfigurationObjectIsNil).CodeU8().Set(eConfigurationObjectIsNil),
	ConfigurationObjectIsNotValid:             dic.NewError(cConfigurationObjectIsNotValid, "тип").CodeU8().Set(eConfigurationObjectIsNotValid),
	ConfigurationObjectIsNotAddress:           dic.NewError(cConfigurationObjectIsNotAddress, "тип").CodeU8().Set(eConfigurationObjectIsNotAddress),
	ConfigurationObjectCopy:                   dic.NewError(cConfigurationObjectCopy, "тип", "ошибка").CodeU8().Set(eConfigurationObjectCopy),
	ConfigurationCallbackAlreadyRegistered:    dic.NewError(cConfigurationCallbackAlreadyRegistered, "тип", "функция").CodeU8().Set(eConfigurationCallbackAlreadyRegistered),
	ConfigurationCallbackSubscriptionNotFound: dic.NewError(cConfigurationCallbackSubscriptionNotFound, "тип", "функция").CodeU8().Set(eConfigurationCallbackSubscriptionNotFound),
}

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors
	// ApplicationPanicException Выполнение приложения прервано паникой:\n...\n...
	ApplicationPanicException dic.IError
	// ApplicationUnknownError Неизвестная ошибка ...
	ApplicationUnknownError dic.IError
	// ApplicationHelpDisplayed Отображение помощи ...
	ApplicationHelpDisplayed dic.IError
	// ApplicationVersion Версии приложения содержит ошибку: ...
	ApplicationVersion dic.IError
	// ApplicationMainFuncNotFound Не определена основная функция приложения.
	ApplicationMainFuncNotFound dic.IError
	// ApplicationMainFuncAlreadyRegistered Основная функция приложения уже зарегистрирована.
	ApplicationMainFuncAlreadyRegistered dic.IError
	// ApplicationRegistrationUnknownObject Регистрация не известного компонента, объекта или модуля: ...
	ApplicationRegistrationUnknownObject dic.IError
	// ComponentIsNull В качестве объекта компоненты передан nil.
	ComponentIsNull dic.IError
	// ComponentRegistrationProhibited Регистрация компонентов запрещена. Компонента ... не зарегистрирована.
	ComponentRegistrationProhibited dic.IError
	// ComponentRegistrationError Регистрация компоненты ... завершилась ошибкой: ...
	ComponentRegistrationError dic.IError
	// ComponentPreferencesCallBeforeCompleting Опрос настроек компонентов вызван до завершения регистрации компонентов.
	ComponentPreferencesCallBeforeCompleting dic.IError
	// ComponentPanicException Выполнение компоненты ... прервано паникой:\n...\n...
	ComponentPanicException dic.IError
	// ComponentRunlevelError Уровень запуска (runlevel), для компоненты ... указан ..., необходимо указать уровень равный 0, либо в интервале от 10 до 65534 включительно.
	ComponentRunlevelError dic.IError
	// ComponentRulesError Правила ... для компоненты ... содержат ошибку: ...
	ComponentRulesError dic.IError
	// RunlevelCantLessCurrentLevel Новый уровень работы приложения (...) не может быть меньше текущего уровня работы приложения (...).
	RunlevelCantLessCurrentLevel dic.IError
	// InitLogging Критическая ошибка в модуле менеджера логирования: ...
	InitLogging dic.IError
	// ComponentConflict Компонента ... конфликтует с компонентой ...
	ComponentConflict dic.IError
	// ComponentRequires Компонента ... имеет не удовлетворённую зависимость ...
	ComponentRequires dic.IError
	// ComponentInitiateTimeout Превышено время ожидание выполнения функции Initiate() компоненты ...
	ComponentInitiateTimeout dic.IError
	// ComponentInitiateExecution Выполнение функции Initiate() компоненты ... завершено с ошибкой: ...
	ComponentInitiateExecution dic.IError
	// ComponentInitiatePanicException Выполнение функции Initiate() компоненты ... прервано паникой:\n...\n...
	ComponentInitiatePanicException dic.IError
	// ComponentDoExecution Выполнение функции Do() компоненты ... завершено с ошибкой: ...
	ComponentDoExecution dic.IError
	// ComponentDoPanicException Выполнение функции Do() компоненты ... прервано паникой:\n...\n...
	ComponentDoPanicException dic.IError
	// ComponentDoUnknownError Выполнение функций Do() завершилось ошибкой: ...
	ComponentDoUnknownError dic.IError
	// ComponentFinalizeExecution Выполнение функции Finalize() компоненты ... завершено с ошибкой: ...
	ComponentFinalizeExecution dic.IError
	// ComponentFinalizePanicException Выполнение функции Finalize() компоненты ... прервано паникой:\n...\n...
	ComponentFinalizePanicException dic.IError
	// ComponentFinalizeUnknownError Выполнение функций Finalize() завершилось ошибкой: ...
	ComponentFinalizeUnknownError dic.IError
	// ComponentFinalizeWarning Выполнение функций Finalize() компоненты ..., длится дольше отведённого времени (...).
	ComponentFinalizeWarning dic.IError
	// RunlevelSubscribeUnsubscribeNilFunction Передана nil функция, подписка или отписка nil функции не возможна.
	RunlevelSubscribeUnsubscribeNilFunction dic.IError
	// RunlevelAlreadySubscribedFunction Функция ... уже подписана на получение событий изменения уровня работы приложения.
	RunlevelAlreadySubscribedFunction dic.IError
	// RunlevelSubscriptionNotFound Не найдена подписка функции ... на события изменения уровня работы приложения.
	RunlevelSubscriptionNotFound dic.IError
	// RunlevelSubscriptionPanicException Вызов функции подписчика на событие изменения уровня работы приложения, прервано паникой:\n...\n....
	RunlevelSubscriptionPanicException dic.IError
	// CommandLineArgumentRequired Требуется указать обязательную команду, аргумент или флаг командной строки: ...
	CommandLineArgumentRequired dic.IError
	// CommandLineArgumentUnknown Неизвестная команда, аргумент или флаг командной строки: ...
	CommandLineArgumentUnknown dic.IError
	// CommandLineArgumentNotCorrect Не верное значение или тип аргумента, флага или параметра: ...
	CommandLineArgumentNotCorrect dic.IError
	// CommandLineRequiredFlag Не указан один или несколько обязательных флагов: ...
	CommandLineRequiredFlag dic.IError
	// CommandLineUnexpectedError Не предвиденная ошибка библиотеки командного интерфейса приложения: ...; ...
	CommandLineUnexpectedError dic.IError
	// ConfigurationBootstrap Ошибка начально bootstrap конфигурации приложения: ...
	ConfigurationBootstrap dic.IError
	// GetCurrentUser Не удалось загрузить данные о текущем пользователе операционной системы: ...
	GetCurrentUser dic.IError
	// CantChangeWorkDirectory Не удалось сменить рабочую директорию приложения: ...
	CantChangeWorkDirectory dic.IError
	// PidExistsAnotherProcessOfApplication Существует один или несколько работающих процессов приложения, измените PID файл или остановите экземпляры приложения, PID: ...
	PidExistsAnotherProcessOfApplication dic.IError
	// PidFileError Ошибка работы с PID файлом ...: ...
	PidFileError dic.IError
	// DatabusRecursivePointer Не возможно определить тип рекурсивного указателя: ...
	DatabusRecursivePointer dic.IError
	// DatabusPanicException Работа с подпиской потребителя, в шине данных, прервана паникой:\n...\n...
	DatabusPanicException dic.IError
	// DatabusSubscribeNotFound Потребитель данных ... не был подписан на шину данных.
	DatabusSubscribeNotFound dic.IError
	// DatabusInternalError Внутренняя ошибка шины данных: ...
	DatabusInternalError dic.IError
	// DatabusNotSubscribersForType Отсутствуют потребители данных для типа данных: ...
	DatabusNotSubscribersForType dic.IError
	// DatabusObjectIsNil Передан nil объект.
	DatabusObjectIsNil dic.IError
	// ConfigurationApplicationProhibited Регистрация объектов конфигурации на текущем уровне работы приложения запрещена. Конфигурация ... не зарегистрирована.
	ConfigurationApplicationProhibited dic.IError
	// ConfigurationApplicationObject Объект конфигурации приложения содержит ошибку: ...
	ConfigurationApplicationObject dic.IError
	// ConfigurationApplicationPanic Непредвиденная ошибка при регистрации объекта конфигурации.\nПаника: ...\n...
	ConfigurationApplicationPanic dic.IError
	// ConfigurationFileNotFound Указанного файла конфигурации ... не существует: ...
	ConfigurationFileNotFound dic.IError
	// ConfigurationPermissionDenied Отсутствует доступ к файлу конфигурации ..., ошибка: ...
	ConfigurationPermissionDenied dic.IError
	// ConfigurationUnexpectedMistakeFileAccess Неожиданная ошибка доступа к файлу конфигурации ...: ...
	ConfigurationUnexpectedMistakeFileAccess dic.IError
	// ConfigurationFileIsDirectory В качестве файла конфигурации указана директория: ...
	ConfigurationFileIsDirectory dic.IError
	// ConfigurationFileReadingError Чтение фала конфигурации ... прервано ошибкой: ...
	ConfigurationFileReadingError dic.IError
	// ConfigurationSetDefault Установка значений по умолчанию, для переменных конфигурации, прервана ошибкой: ...
	ConfigurationSetDefault dic.IError
	// ConfigurationSetDefaultValue Установка значения по умолчанию ..., для переменной конфигурации ..., прервана ошибкой: ...
	ConfigurationSetDefaultValue dic.IError
	// ConfigurationSetDefaultPanic Непредвиденная ошибка, при установке значений по умолчанию, объекта конфигурации.\nПаника: ...\n...
	ConfigurationSetDefaultPanic dic.IError
	// ConfigurationObjectNotFound Объект конфигурации с типом ... не найден.
	ConfigurationObjectNotFound dic.IError
	// ConfigurationObjectIsNotStructure Переданный объект ... не является структурой.
	ConfigurationObjectIsNotStructure dic.IError
	// ConfigurationObjectIsNil Переданный объект, является nil объектом.
	ConfigurationObjectIsNil dic.IError
	// ConfigurationObjectIsNotValid Объект конфигурации с типом ... не инициализирован.
	ConfigurationObjectIsNotValid dic.IError
	// ConfigurationObjectIsNotAddress Объект конфигурации с типом ... передан не корректно. Необходимо передать адрес объекта.
	ConfigurationObjectIsNotAddress dic.IError
	// ConfigurationObjectCopy Копирование объекта конфигурации с типом ... прервано ошибкой: ...
	ConfigurationObjectCopy dic.IError
	// ConfigurationCallbackAlreadyRegistered Подписка функции обратного вызова на изменение конфигурации с типом ... для функции ... уже существует.
	ConfigurationCallbackAlreadyRegistered dic.IError
	// ConfigurationCallbackSubscriptionNotFound Подписка функции обратного вызова на изменение конфигурации с типом ... для функции ... не существует.
	ConfigurationCallbackSubscriptionNotFound dic.IError
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
