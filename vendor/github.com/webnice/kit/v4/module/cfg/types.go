package cfg

import (
	"container/list"
	"os"
	"os/user"
	"regexp"
	"sync"
	"time"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitModuleBus "github.com/webnice/kit/v4/module/bus"
	kitModuleCfgCli "github.com/webnice/kit/v4/module/cfg/cli"
	kitModuleLog "github.com/webnice/kit/v4/module/log"
	kitModuleServer "github.com/webnice/kit/v4/module/server"
	kitModuleUuid "github.com/webnice/kit/v4/module/uuid"
	kitTypes "github.com/webnice/kit/v4/types"

	"github.com/Masterminds/semver/v3"
)

var (
	singleton                       *impl  // Конфигурация всегда в единичном экземпляре.
	defaultWorkingDirectoryOriginal string // Рабочая директория, первоначальное, скорректированное значение.
	defaultTempDirectoryOriginal    string // Директория временных файлов, первоначальное, скорректированное значение.
	defaultCacheDirectoryOriginal   string // Директория для файлов кеша, первоначальное, скорректированное значение.
	defaultConfigDirectoryOriginal  string // Директория для конфигураций, первоначальное, скорректированное значение.
	defaultLogDirectoryOriginal     string // Директория для файлов журнала приложения, первоначальное, скорректированное значение.
)

// Regexp коррекции печати конфигурации приложения.
var rexUnionStructureHeader = regexp.MustCompile(`(?im)\*struct\s+({[^}]+})`)

// Объект сущности пакета.
type impl struct {
	bootstrapConfiguration    *kitTypes.BootstrapConfiguration           // Минимальная, стартовая конфигурация приложения.
	loadableConfiguration     *kitTypes.BootstrapConfiguration           // Загружаемая из файла, временная конфигурация.
	isForkWorker              bool                                       // Режим forkWorker.
	command                   []string                                   // Команда с которой было запущено приложение.
	error                     []error                                    // Список ошибок возникших при выполнении компонентов приложения.
	main                      *mainApplication                           // Основная функция приложения.
	comp                      *mainComponent                             // Регистр компонентов приложения.
	conf                      *mainConfiguration                         // Регистр объектов конфигураций.
	runLevel                  uint16                                     // Текущий уровень работы приложения.
	mapRunLevel               []uint16                                   // Карта переключения уровней работы приложения, построенная по уровням компонентов.
	runLevelChangeChan        chan *runLevelUp                           // Канал изменения уровня работы приложения.
	runLevelSubscribers       *list.List                                 // Список подписчиков на события изменения уровня работы приложения.
	runLevelStopAutoincrement bool                                       // Остановка автоматического переключения уровня работы приложения.
	finalizeWarningTimeout    time.Duration                              // Время отводимое на выполнение функции Finalize() до печали в лог warning сообщения.
	cacheForkWorker           *kitTypes.BootstrapConfigurationForkWorker // Параметры работы в режиме fork worker - кеш.
	user                      *user.User                                 // Текущий пользователь операционной системы.
	rawWriter                 kitTypes.SyncWriter                        // Интерфейс вывода потоковых сообщений.
	noPidFile                 bool                                       // Запрет работы плагина pidfile.
	noDbMigration             bool                                       // Запрет применения миграций базы данных.

	// Выделенные в отдельные сущности интерфейсы.

	bus     kitModuleBus.Interface    // Объект интерфейса Databus.
	version Version                   // Объект интерфейса Version.
	essence Essence                   // Объект интерфейса Essence.
	rec     Recorder                  // Объект интерфейса Recorder.
	logger  kitModuleLog.Logger       // Объект интерфейса Logger.
	uuid    kitModuleUuid.Interface   // Объект интерфейса UUID.
	cli     kitModuleCfgCli.Interface // Объект интерфейса CLI.
	srv     kitModuleServer.Interface // Объект интерфейса Server.
	ans     kitModuleAns.Interface    // Объект интерфейса Ans.
}

// Объект сути сущности, интерфейс Essence.
type gist struct {
	parent *impl // Адрес объекта основной сущности, интерфейс Interface.
}

// Объект реализующий интерфейс kitTypes.SyncWriter.
type syncWriter struct {
	parent *impl    // Адрес объекта основной сущности, интерфейс Interface.
	wrh    *os.File // Писатель.
}

// Структура описания основной функции приложения и версии приложения.
type mainApplication struct {
	Version *semver.Version // Версия приложения.
	Fn      kitTypes.MainFn // Функция входа, выполнения и завершения приложения с кодом ошибки.
	FnMutex *sync.Mutex     // Обеспечение запуска только одной копии приложения.
}

// Структура регистра компонентов приложения.
type mainComponent struct {
	IsClose        bool                      // Флаг закрытия регистрации компонентов приложения.
	Registered     []*kitTypes.ComponentInfo // Временный список зарегистрированных компонентов.
	Groups         []*kitTypes.CommandGroup  // Описания группы команд компонентов приложения.
	Component      []*kitTypes.ComponentInfo // Список компонентов приложения, отсортированных в очерёдности запуска.
	ComponentMutex *sync.Mutex               // Защита данных.
}
