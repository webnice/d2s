package cfg

import (
	"io"
	"os/user"
	"reflect"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitModuleBus "github.com/webnice/kit/v4/module/bus"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitModuleServer "github.com/webnice/kit/v4/module/server"
	kitModuleUuid "github.com/webnice/kit/v4/module/uuid"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Interface Интерфейс пакета.
type Interface interface {
	// ИНТЕРФЕЙСЫ.

	// Answer Интерфейс библиотеки функций для формирования ответа на HTTP запрос к серверу.
	Answer() kitModuleAns.Interface

	// Bus Интерфейс шины данных приложения.
	Bus() kitModuleBus.Interface

	// Gist Интерфейс к служебным методам конфигурации приложения.
	Gist() Essence

	// Version Интерфейс к методам получения версии приложения.
	Version() Version

	// UUID Интерфейс к методам генерации и работы с UUID идентификаторами.
	UUID() kitModuleUuid.Interface

	// RawWriter Интерфейс вывода потоковых сообщений.
	RawWriter() kitTypes.SyncWriter

	// Server Интерфейс менеджера управления сервером.
	Server() kitModuleServer.Interface

	// МЕНЕДЖЕР ЛОГИРОВАНИЯ.

	// Log Интерфейс к методам логирования.
	Log() kitTypes.Logger

	// Loglevel Текущий уровень логирования приложения.
	Loglevel() kitModuleLogLevel.Level

	// МЕТОДЫ.

	// AppName Название приложения.
	AppName() string

	// Command Функция возвращает команду с которой было запущено приложение.
	Command() (ret string)

	// CommandFull Функция возвращает команду с которой было запущено приложение.
	CommandFull() []string

	// Debug Функция возвращает текуще значение режима отладки приложения.
	Debug() (ret bool)

	// ForkWorker Параметры работы в режиме fork worker.
	ForkWorker() *kitTypes.BootstrapConfigurationForkWorker

	// IsForkWorker Возвращает истину для режима запуска приложения в качестве подчинённого процесса основного приложения.
	// В этом режиме приложению нельзя использовать STDIN/STDOUT/STDERR, эти потоки используются для обмена данными
	// с управляющим приложением.
	// Так же, этот режим устанавливается в случае запуска приложения как daemon процесса, в этом режиме запуска, тоже
	// нельзя использовать STDIN/STDOUT/STDERR, так как режим daemon, для приложения, возможен только после полного
	// закрытия этих потоков.
	IsForkWorker() (ret bool)

	// Runlevel Возвращает текущее значение уровня работы приложения.
	Runlevel() uint16

	// RunlevelMap Возвращает карту, описывающую план переключения уровней работы приложения.
	RunlevelMap() []uint16

	// RunlevelSubscribe Подписка на события изменения уровня работы приложения.
	RunlevelSubscribe(fn RunlevelSubscriberFn) (err error)

	// RunlevelUnsubscribe Отписка от событий изменения уровня работы приложения.
	RunlevelUnsubscribe(fn RunlevelSubscriberFn) (err error)

	// RunlevelAutoincrement Режим работы автоматического увеличения уровня работы приложения.
	RunlevelAutoincrement() (ret bool)

	// Targetlevel Возвращает текущее значение целевого уровня работы приложения.
	Targetlevel() uint16

	// User Текущий пользователь операционной системы.
	User() *user.User

	// AbsolutePath Функция преобразует относительный путь в абсолютный путь к директории или файлу.
	// Учитывается символ '~' обозначающий домашнюю директорию текущего пользователя.
	AbsolutePath(pth string) (ret string)

	// DirectoryHome Значение домашней директории приложения.
	DirectoryHome() string

	// DirectoryWorking Значение рабочей директории приложения.
	DirectoryWorking() string

	// DirectoryWorkingChdir Функция выполняет переход в рабочую директорию приложения.
	// Возвращается ошибка операционной системы, с причиной не возможности перейти в рабочую директорию.
	DirectoryWorkingChdir() (err error)

	// DirectoryTemp Значение директории для временных файлов.
	DirectoryTemp() string

	// DirectoryCache Значение директории для файлов кеша.
	DirectoryCache() string

	// DirectoryConfig Значение директории для подключаемых или дополнительных конфигураций приложения.
	DirectoryConfig() string

	// DirectoryLog Значение директории для файлов журнала.
	DirectoryLog() string

	// FileConfig Значение пути и имени конфигурационного файла приложения.
	FileConfig() string

	// FilePid Значение пути и имени PID файла приложения.
	FilePid() string

	// FileState Значение пути и имени файла хранения состояния приложения.
	FileState() string

	// FileSocket Значение пути и имени сокет файла коммуникаций с приложением.
	FileSocket() string

	// FileLogWithName Создание или открытие на запись файла журнала, в директории для журналов, с указанием
	// названия файла журнала. Возвращается писатель, который необходимо закрыть, после окончания записи данных в журнал.
	FileLogWithName(logName string) (ret io.WriteCloser, err error)

	// КОНФИГУРАЦИЯ.

	// ConfigurationUnionSprintf Печать объединённой конфигурации приложения в строку.
	ConfigurationUnionSprintf() (ret string)

	// ConfigurationByType Возвращает объект конфигурации соответствующий указанному типу объекта.
	// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
	ConfigurationByType(t reflect.Type) (ret any, err error)

	// ConfigurationByTypeName Возвращает объект конфигурации соответствующий указанному названию типа объекта.
	// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
	ConfigurationByTypeName(typeName string) (ret any, err error)

	// ConfigurationByObject Возвращает объект конфигурации соответствующий типу переданного объекта, сам переданный
	// объект никак не изменяется, он служит только для определения типа данных.
	// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
	ConfigurationByObject(o any) (ret any, err error)

	// ConfigurationCopyByObject Если существует конфигурация с типом данных идентичным переданному объекту,
	// тогда данные конфигурации копируются в переданный объект.
	// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
	// Объект должен передаваться по адресу, иначе его заполнение не возможно и будет возвращена ошибка.
	ConfigurationCopyByObject(o any) (err error)

	// ОШИБКИ

	// Errors Справочник ошибок.
	Errors() *Error
}
