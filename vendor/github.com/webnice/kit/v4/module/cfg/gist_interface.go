package cfg

import (
	"bytes"
	"time"

	kitModuleCfgCli "github.com/webnice/kit/v4/module/cfg/cli"
	kitModuleLog "github.com/webnice/kit/v4/module/log"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Essence Служебный публичный интерфейс.
type Essence interface {
	// ИНТЕРФЕЙСЫ.

	// Cfg Интерфейс конфигурации.
	Cfg() Interface

	// Logger Интерфейс менеджера логирования.
	Logger() kitModuleLog.Logger

	// CLI Интерфейс к методам работы с параметрами командной строки и переменных окружения.
	CLI() kitModuleCfgCli.Interface

	// МЕТОДЫ.

	// App Точка запуска приложения.
	App()

	// AppName Установка значения названия приложения.
	AppName(name string) Essence

	// Version Присвоение версии и номера сборки приложения в семантике "Semantic Versioning 2.0.0".
	Version(version string, build string) Essence

	// Debug Присвоение нового значения режима отладки приложения.
	Debug(debug bool) Essence

	// NoPidFileSet Установка значения флага запрета работы плагина pidfile.
	// Используется для режима работы приложения в качестве утилиты командной строки.
	NoPidFileSet(noPidFile bool) Essence

	// NoPidFileGet Получение значения флага запрета работы плагина pidfile.
	NoPidFileGet() bool

	// NoDbMigrationSet Установка значения флага запрета применения миграций базы данных.
	// Используется для режима работы приложения в качестве утилиты командной строки.
	NoDbMigrationSet(noDbMigration bool) Essence

	// NoDbMigrationGet Получение значения флага запрета применения миграций базы данных.
	NoDbMigrationGet() bool

	// UpdateBootstrapConfigurationPathValue Коррекция значений стартовой конфигурация приложения, если они были загружены
	// из командной строки или переменной окружения, то есть не равны значениям по умолчанию
	UpdateBootstrapConfigurationPathValue() Essence

	// CommandFull Присвоение значения команды с которой было запущено приложение.
	CommandFull(cmd []string) Essence

	// ForkWorkerMode Присвоение нового значения режима запуска в качестве подчинённого процесса основного приложения.
	ForkWorkerMode(isForkWorker bool) Essence

	// ErrorAppend Добавление ошибки в стек ошибок приложения.
	ErrorAppend(err error) Essence

	// IsError Возвращает истину, если есть ошибки приложения.
	IsError() bool

	// Runlevel Установка или увеличение значения уровня работы приложения.
	// Поведение функции в зависимости от значения аргумента runlevel:
	// * runlevel>0 - будет выполнена установка нового значения, новое значение не может быть меньше текущего значения.
	// * runlevel=0 - значение Runlevel будет увеличено на единицу.
	Runlevel(runlevel uint16) Essence

	// RunlevelNext Переключение значения уровня работы приложения на следующий уровень.
	// От 01 до 10 переключается на единицу, последовательно.
	// От 10 до 65535 переключается согласно карте уровней, если текущий уровень не входит в карту - последовательно.
	// На уровне 65535 переключение останавливается, вызов функции не меняет уровень.
	RunlevelNext() Essence

	// RunlevelNextAsync Асинхронное переключение значения уровня работы приложения на следующий уровень.
	RunlevelNextAsync()

	// RunlevelDefault Уровень выполнения компонентов приложения по умолчанию.
	RunlevelDefault() (ret uint16)

	// RunlevelAutoincrementStop Отключение режима автоматического увеличения уровня работы приложения.
	RunlevelAutoincrementStop() Essence

	// RunlevelAutoincrementStart Включение режима автоматического увеличения уровня работы приложения.
	RunlevelAutoincrementStart() Essence

	// RunlevelExit Переключение значения уровня работы приложения на уровень завершения работы приложения.
	RunlevelExit() Essence

	// RunlevelExitAsync Переключение значения уровня работы приложения на уровень завершения работы приложения.
	RunlevelExitAsync()

	// Targetlevel Установка значения целевого уровня работы приложения.
	Targetlevel(tl uint16) Essence

	// Registration Регистрация разных объектов приложения.
	Registration(name string, obj any) (err error)

	// КОМПОНЕНТЫ.

	// ComponentName Получение уникального имени пакета компоненты.
	ComponentName(obj any) (ret string)

	// ComponentNames Возвращает список зарегистрированных компонентов.
	ComponentNames() []string

	// ComponentPreferences Функция-менеджер загрузки и обработки настроек компонентов.
	ComponentPreferences(fn kitTypes.ComponentPreferencesFn) (code uint8, err error)

	// ComponentCheckConflict Проверка конфликтов между всеми зарегистрированными компонентами.
	ComponentCheckConflict(fn kitTypes.ComponentConflictFn) (code uint8, err error)

	// ComponentRequiresCheck Проверка зависимостей между всеми зарегистрированными компонентами.
	ComponentRequiresCheck(fn kitTypes.ComponentRequiresFn) (code uint8, err error)

	// ComponentSort Сортировка зарегистрированных компонентов в соответствии с настройками (before) и (after).
	ComponentSort(fn kitTypes.ComponentSortFn) (code uint8, err error)

	// ComponentMapRunlevel Построение шагов переключения уровня выполнения приложения (runlevel).
	ComponentMapRunlevel(targetlevel uint16) (code uint8, err error)

	// ComponentInitiate Вызов функции Initiate у всех зарегистрированных компонентов в прямом порядке.
	ComponentInitiate(fn kitTypes.ComponentInitiateFn) (code uint8, err error)

	// ComponentDo Вызов функции Do у всех зарегистрированных компонентов в прямом порядке для указанного уровня приложения.
	ComponentDo(runlevel uint16, fn kitTypes.ComponentDoFn) (code uint8, err error)

	// ComponentFinalizeWarningTimeout Возвращает время отводимое на выполнение функции Finalize(), до печати в лог
	// сообщения о долгой работе функции.
	ComponentFinalizeWarningTimeout() (ret time.Duration)

	// ComponentFinalize Вызов функции Finalize у всех зарегистрированных компонентов в обратном порядке.
	ComponentFinalize(fn kitTypes.ComponentFinalizeFn) (code uint8, err error)

	// ComponentCommandRegister Регистрация команды и группы команд компоненты.
	ComponentCommandRegister(cc kitTypes.ComponentCommand)

	// ComponentFlagRegister Регистрация глобального флага компоненты приложения.
	ComponentFlagRegister(cf kitTypes.ComponentFlag)

	// ДИРЕКТОРИИ И ФАЙЛЫ ОСНОВНОЙ КОНФИГУРАЦИИ.

	// DirectoryHome Установка значения домашней директория приложения.
	DirectoryHome(ahd string) Essence

	// DirectoryWorking Установка значения рабочей директории приложения.
	DirectoryWorking(awd string) Essence

	// DirectoryTemp Установка значения директории для временных файлов.
	DirectoryTemp(atd string) Essence

	// DirectoryCache Установка значения директории для файлов кеша.
	DirectoryCache(acd string) Essence

	// DirectoryConfig Установка значения директории для подключаемых или дополнительных конфигураций приложения.
	DirectoryConfig(afd string) Essence

	// FileConfig Установка значения пути и имени конфигурационного файла приложения.
	FileConfig(cfn string) Essence

	// FilePid Установка значения пути и имени PID файла приложения.
	FilePid(pfn string) Essence

	// FileState Установка значения пути и имени файла хранения состояния приложения.
	FileState(stf string) Essence

	// FileSocket Установка значения пути и имени сокет файла коммуникаций с приложением.
	FileSocket(sof string) Essence

	// LogLevel Установка значения уровня логирования приложения.
	LogLevel(ll kitModuleLogLevel.Level) Essence

	// КОНФИГУРАЦИЯ.

	// ConfigurationRegistration Регистрация объекта конфигурации, являющегося частью общей конфигурации приложения.
	// Регистрация доступна только на уровне работы приложения 0.
	// Объект конфигурации должен передаваться в качестве адреса.
	// Поля объекта конфигурации должны состоять из простых и/или сериализуемых типов данных и быть экспортируемыми.
	// Поля объекта могут содержать теги, которые определяют внедрение конфигурации в конфигурацию приложения.
	// Вместе с объектом конфигурации можно передать функцию обратного вызова, она будет вызвана при изменении данных
	// конфигурации, например при перезагрузке файла конфигурации или иных реализациях динамического изменения
	// значений конфигурации.
	// Теги:
	//   - description ---- Описание поля, публикуется в YAML файле, при создании примера конфигурации,
	//     подробности в компоненте "configuration".
	//     Если указано значение "-", тогда описание не публикуется.
	//   - default-value -- Значение поля по умолчанию, присваивается после чтения конфигурационного файла,
	//     а так же, публикуется в YAML файле, при создании примера конфигурации.
	//   - yaml ----------- Тег для библиотеки YAML, если указано значение "-", тогда поле пропускается.
	// Возвращаемый результат:
	// - истина - в случае успешного внедрения конфигурации в конфигурацию приложения;
	// - ложь   - в случае возникновения ошибки при внедрении конфигурации. Сама ошибка публикуется в список ошибок
	//            приложения.
	ConfigurationRegistration(c any, callback ...kitTypes.Callbacker) (isOk bool)

	// ConfigurationCallbackSubscribe Подписка функции обратного вызова на событие изменения данных сегмента
	// конфигурации. Функция будет вызвана при изменении данных конфигурации, например при перезагрузке файла
	// конфигурации или иных реализациях динамического изменения значений конфигурации.
	ConfigurationCallbackSubscribe(c any, callback kitTypes.Callbacker) (err error)

	// ConfigurationCallbackUnsubscribe Отписка функции обратного вызова на событие изменения данных сегмента
	// конфигурации.
	ConfigurationCallbackUnsubscribe(c any, callback kitTypes.Callbacker) (err error)

	// ConfigurationLoad Загрузка конфигурационного файла.
	ConfigurationLoad(buf *bytes.Buffer) (err error)

	// AbsolutePathAndUpdate Функция преобразует относительный путь в абсолютный путь,
	// проверяет равно ли новое значение старому, если значение изменилось, тогда обновляет значение в указанной
	// переменной и возвращает результирующее значение.
	AbsolutePathAndUpdate(dir *string) (ret string)
}
