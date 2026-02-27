package cli

import (
	"bytes"

	kitModuleCfgCliKong "github.com/webnice/kit/v4/module/cfg/cli/kong"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Interface Интерфейс пакета.
type Interface interface {
	// Constant Установка названия переменных окружения, используемых отдельно от структуры конфигурации.
	Constant(env ConstantEnvironmentName) Interface

	// Bootstrap Первоначальная инициализация минимальной конфигурации приложения.
	Bootstrap(b *kitTypes.BootstrapDefaultValue) (err error)

	// RegisterCommand Регистрации динамических команд приложения.
	RegisterCommand(cmd *Command)

	// RegisterFlag Регистрация динамических глобальных флагов приложения.
	RegisterFlag(flg *Flag)

	// Init Инициализация командного интерфейса и загрузка переменных окружения.
	Init() (help *bytes.Buffer, description string, err error)

	// Command Команда приложения, первая часть.
	Command() (ret string)

	// CommandFull Полная команда приложения.
	CommandFull() []string

	// ОШИБКИ.

	// Errors Справочник всех ошибок пакета.
	Errors() *Error
}

// Объект сущности пакета.
type impl struct {
	env                   ConstantEnvironmentName          // Названия специальных переменных окружения.
	bootstrap             *kitTypes.BootstrapConfiguration // Базовая конфигурация приложения.
	bootstrapDefaultValue *kitTypes.BootstrapDefaultValue  // Минимальные начальные данные необходимые для работы приложения.
	applicationName       string                           // Название приложения.
	applicationCommand    []string                         // Полученная команда приложения.
	k2g                   *kitModuleCfgCliKong.Kong        // Объект библиотеки kong.
	k2gOut                *bytes.Buffer                    // kong stdout.
	k2gErr                *bytes.Buffer                    // kong stderr.
	k2gIsExit             bool                             // Флаг устанавливается в истина, если cli вызвала функцию завершения.
	k2gExitCode           int                              // Код завершения, возвращённый cli.
	ctx                   *kitModuleCfgCliKong.Context     // Контекст библиотеки kong.
	command               []*Command                       // Динамические команды приложения.
	flag                  []*Flag                          // Динамические глобальные флаги приложения.
}

// Command Структура регистрации динамических команд приложения.
type Command struct {
	// Command Название команды. Пустые команды не создаются.
	Command string

	// Description Описание команды, отображается в помощи пользователю.
	Description string

	// GroupKey Идентификатор группы команд, ключ связи описания группы с командой, входящей в группу.
	GroupKey string

	// IsDefault Команда по умолчанию.
	IsDefault bool

	// IsHidden Команда скрыта, не отображается в помощи.
	IsHidden bool

	// Value Ссылка на структуру значений с мета информацией, в неё же будут загружены
	// указанное в CLI или ENV, значения.
	Value any
}

// Flag Структура регистрации динамических глобальных флагов приложения.
type Flag struct {
	// ShortKey Короткий, односимвольный синоним флага, может быть пустым.
	ShortKey rune

	// Flag Полное наименование флага, если пустой, флаг игнорируется.
	Flag string

	// Description Помощь для пользователя, описывающая назначение флага.
	Description string

	// Environment Наименование переменной окружения из которой может быть взято значение флага.
	Environment string

	// Placeholder Значение флага, которое будет указано при отображении в помощи.
	Placeholder string

	// Обязательный флаг.
	IsRequired bool

	// IsHidden Флаг скрыт, не отображается в помощи.
	IsHidden bool

	// Value Ссылка на переменную, в которую будет загружено указанное в CLI или ENV, значение.
	Value any
}
