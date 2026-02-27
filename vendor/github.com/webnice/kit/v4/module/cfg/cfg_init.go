package cfg

import (
	"container/list"
	"sync"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitModuleBus "github.com/webnice/kit/v4/module/bus"
	kitModuleCfgCli "github.com/webnice/kit/v4/module/cfg/cli"
	kitModuleLog "github.com/webnice/kit/v4/module/log"
	kitModuleServer "github.com/webnice/kit/v4/module/server"
	kitModuleUuid "github.com/webnice/kit/v4/module/uuid"
	kitTypes "github.com/webnice/kit/v4/types"

	"github.com/Masterminds/semver/v3"
)

func init() {
	singleton = &impl{
		// Минимальная, стартовая конфигурация приложения.
		bootstrapConfiguration: &kitTypes.BootstrapConfiguration{ApplicationDebug: initApplicationDebug()},
		// Загружаемая из файла, временная конфигурация.
		loadableConfiguration: &kitTypes.BootstrapConfiguration{},
		// Основная функция приложения.
		main: &mainApplication{
			Version: &semver.Version{}, // Исключение возможной ошибки, если версия не установлена.
			Fn:      nil,               // Функция входа, выполнения и завершения приложения с кодом ошибки.
			FnMutex: new(sync.Mutex),   // Обязательная защита от двойного запуска основной функции приложения.
		},
		// Регистр компонентов приложения.
		comp: &mainComponent{
			IsClose:        false,
			ComponentMutex: new(sync.Mutex),
		},
		// Регистр объектов конфигураций.
		conf: &mainConfiguration{
			Items: make([]*configurationItem, 0, 1),
		},
		// Время отводимое на выполнение функции Finalize() до печали в лог warning сообщения.
		finalizeWarningTimeout: defaultFinalizeTimeoutBeforeWarning(),
	}
	// Интерфейс вывода потоковых сообщений.
	singleton.rawWriter = newSyncWriter(singleton)
	// Пользователь операционной системы.
	singleton.user = currentUser(singleton)
	// Добавление интерфейсов используемых модулей.
	singleton.bus = kitModuleBus.New(0, 0, singleton.bootstrapConfiguration.ApplicationDebug)
	singleton.version, singleton.essence = newVersion(singleton), newEssence(singleton)
	singleton.rec, singleton.logger, singleton.uuid =
		newRecorder(singleton), kitModuleLog.New(singleton.rawWriter, singleton.bus), kitModuleUuid.Get()
	// Основная часть конфигурации приложения, интерфейс работы с командной строкой и переменными окружения.
	singleton.initMainConfiguration()
	// Управление переключением уровней работы приложения.
	singleton.runLevelChangeChan, singleton.runLevelSubscribers = runlevelChangeFuncNew(), list.New()
	// Интерфейс сервера.
	singleton.srv = kitModuleServer.New(singleton.rec)
	// Интерфейс ans.Interface.
	singleton.ans = kitModuleAns.New(singleton.rec)
}

// Основная часть конфигурации приложения, интерфейс работы с командной строкой и переменными окружения.
func (cfg *impl) initMainConfiguration() {
	var (
		empty func() string
		err   error
	)

	cfg.cli = kitModuleCfgCli.
		New(cfg.bootstrapConfiguration).
		Constant(defaultConstantEnvironment())
	empty = func() (ret string) { return }
	if err = cfg.cli.Bootstrap(&kitTypes.BootstrapDefaultValue{
		ApplicationTargetlevel: defaultApplicationTargetlevelString, // Целевой уровень выполнения приложения.
		ApplicationDebug:       defaultApplicationDebugString,       // Режим отладки приложения.
		ApplicationName:        defaultApplicationName,              // Название приложения.
		ApplicationDescription: defaultApplicationDescription,       // Описание приложения или сервиса, отображаемое в строке помощи.
		HomeDirectory:          defaultHomeDirectory,                // Домашняя директория пользователя от которого запущено приложение.
		WorkingDirectory:       defaultWorkingDirectory,             // Рабочая директория приложения.
		TempDirectory:          defaultTempDirectory,                // Директория для создания временных файлов (очищаемая системой).
		CacheDirectory:         defaultCacheDirectory,               // Директория для кеша (не очищаемая системой).
		ConfigDirectory:        defaultConfigDirectory,              // Директория конфигурации в домашней директории пользователя.
		LogDirectory:           defaultLogDirectory,                 // Директория для файлов журнала приложения.
		ConfigFile:             empty,                               // Путь и имя конфигурационного файла по умолчанию.
		PidFile:                empty,                               // Путь и имя PID файла по умолчанию.
		StateFile:              empty,                               // Файл состояния.
		SocketFile:             empty,                               // Сокет файл для коммуникации серверной части с клиентской.
		LogLevel:               defaultLogLevel().String,            // Уровень логирования по умолчанию до загрузки конфигурации.
		CommandGroup:           cfg.initGroupOfCommand,              // Список групп команд, который сформируется при регистрации компонентов.
	}); err != nil {
		cfg.error = append(cfg.error, Errors().ConfigurationBootstrap.Bind(err))
	}
	// Регистрация объекта конфигурации для последующего чтения из файла.
	cfg.essence.ConfigurationRegistration(cfg.loadableConfiguration)
}

// Возврат списка групп команд, который сформируется при регистрации компонентов.
func (cfg *impl) initGroupOfCommand() (ret []kitTypes.CommandGroup) {
	var n int

	if cfg.comp == nil || len(cfg.comp.Groups) <= 0 {
		return
	}
	ret = make([]kitTypes.CommandGroup, 0, len(cfg.comp.Groups))
	for n = range cfg.comp.Groups {
		ret = append(ret, kitTypes.CommandGroup{
			Key:         cfg.comp.Groups[n].Key,
			Title:       cfg.comp.Groups[n].Title,
			Description: cfg.comp.Groups[n].Description,
		})
	}

	return
}
