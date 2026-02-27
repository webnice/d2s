package cfg

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"reflect"
	"strings"

	"github.com/webnice/debug"
	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitModuleBus "github.com/webnice/kit/v4/module/bus"
	kitModuleCfgCpy "github.com/webnice/kit/v4/module/cfg/cpy"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitModuleServer "github.com/webnice/kit/v4/module/server"
	kitModuleUuid "github.com/webnice/kit/v4/module/uuid"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Get Получение объекта сущности пакета, возвращается интерфейс пакета.
//
//goland:noinspection GoUnusedExportedFunction
func Get() Interface { return singleton }

// RegistrationMain Регистрация основной функции приложения.
func RegistrationMain(fn kitTypes.MainFn) Interface {
	var fnName string

	singleton.main.FnMutex.Lock()
	if singleton.main.Fn != nil {
		fnName = getFuncFullName(RegistrationMain)
		singleton.error = append(singleton.error, Errors().ApplicationMainFuncAlreadyRegistered.Bind())
		panic(fnName + "\n" + Errors().ApplicationMainFuncAlreadyRegistered.Bind().Error())
		return singleton
	}
	singleton.main.Fn = fn
	singleton.main.FnMutex.Unlock()

	return singleton
}

// Answer Интерфейс библиотеки функций для формирования ответа на HTTP запрос к серверу.
func (cfg *impl) Answer() kitModuleAns.Interface { return cfg.ans }

// Bus Интерфейс шины данных приложения.
func (cfg *impl) Bus() kitModuleBus.Interface { return cfg.bus }

// Gist Интерфейс к служебным методам.
func (cfg *impl) Gist() Essence { return cfg.essence }

// Version Интерфейс к методам получения версии приложения.
func (cfg *impl) Version() Version { return cfg.version }

// UUID Интерфейс к методам генерации и работы с UUID идентификаторами.
func (cfg *impl) UUID() kitModuleUuid.Interface { return cfg.uuid }

// RawWriter Интерфейс вывода потоковых сообщений.
func (cfg *impl) RawWriter() kitTypes.SyncWriter { return cfg.rawWriter }

// Server Интерфейс менеджера управления сервером.
func (cfg *impl) Server() kitModuleServer.Interface { return cfg.srv }

// Log Интерфейс к методам логирования.
func (cfg *impl) Log() kitTypes.Logger { return cfg.rec }

// Loglevel Значение уровня логирования приложения.
func (cfg *impl) Loglevel() kitModuleLogLevel.Level {
	return kitModuleLogLevel.ParseLevelInt64(int64(cfg.bootstrapConfiguration.LogLevel))
}

// AppName Название приложения.
func (cfg *impl) AppName() string { return cfg.bootstrapConfiguration.ApplicationName }

// Command Функция возвращает команду с которой было запущено приложение.
func (cfg *impl) Command() (ret string) {
	if len(cfg.command) > 0 {
		ret = cfg.command[0]
	}

	return
}

// CommandFull Функция возвращает команду с которой было запущено приложение.
func (cfg *impl) CommandFull() []string { return cfg.command }

// Debug Функция возвращает текуще значение режима отладки приложения.
func (cfg *impl) Debug() (ret bool) { return cfg.bootstrapConfiguration.ApplicationDebug }

// ForkWorker Параметры работы в режиме fork worker.
func (cfg *impl) ForkWorker() *kitTypes.BootstrapConfigurationForkWorker {
	if cfg.cacheForkWorker == nil {
		cfg.cacheForkWorker = &kitTypes.BootstrapConfigurationForkWorker{
			Master:    cfg.bootstrapConfiguration.ApplicationForkWorkerMaster,
			Component: cfg.bootstrapConfiguration.ApplicationForkWorkerComponent,
			Target:    cfg.bootstrapConfiguration.ApplicationForkWorkerTarget,
		}
	}

	return cfg.cacheForkWorker
}

// IsForkWorker Возвращает истину для режима запуска приложения в качестве подчинённого процесса основного приложения.
func (cfg *impl) IsForkWorker() (ret bool) { return cfg.isForkWorker }

// Errors Справочник ошибок.
func (cfg *impl) Errors() *Error { return Errors() }

// Targetlevel Возвращает текущее значение целевого уровня работы приложения.
func (cfg *impl) Targetlevel() uint16 { return cfg.bootstrapConfiguration.ApplicationTargetlevel }

// User Текущий пользователь операционной системы.
func (cfg *impl) User() *user.User { return cfg.user }

// AbsolutePath Функция преобразует относительный путь в абсолютный путь к директории или файлу.
// Учитывается символ '~' обозначающий домашнюю директорию текущего пользователя.
func (cfg *impl) AbsolutePath(pth string) (ret string) { return AbsolutePath(pth) }

// AbsolutePathAndUpdate Функция преобразует относительный путь в абсолютный путь,
// проверяет равно ли новое значение старому, если значение изменилось, тогда обновляет значение в указанной
// переменной и возвращает результирующее значение.
func (cfg *impl) AbsolutePathAndUpdate(dir *string) (ret string) {
	if dir == nil {
		return
	}
	if ret = cfg.AbsolutePath(*dir); *dir != ret {
		*dir = ret
	}

	return
}

// DirectoryHome Значение домашней директории приложения.
func (cfg *impl) DirectoryHome() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.HomeDirectory)
}

// DirectoryWorking Значение рабочей директории приложения.
func (cfg *impl) DirectoryWorking() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.WorkingDirectory)
}

// DirectoryWorkingChdir Функция выполняет переход в рабочую директорию приложения.
// Возвращается ошибка операционной системы, с причиной не возможности перейти в рабочую директорию.
func (cfg *impl) DirectoryWorkingChdir() (err error) {
	var (
		awd string
		cwd string
	)

	awd = cfg.DirectoryWorking()
	cwd, _ = os.Getwd()
	// Рабочая директория из конфигурации совпадает с текущей рабочей директорией, смена директории не требуется.
	if strings.EqualFold(cwd, awd) {
		return
	}
	// Системный вызов смены рабочей директории.
	if err = os.Chdir(awd); err != nil {
		err = fmt.Errorf("syscall chdir(%q) error: %s", awd, err)
		return
	}

	return
}

// DirectoryTemp Значение директории для временных файлов.
func (cfg *impl) DirectoryTemp() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.TempDirectory)
}

// DirectoryCache Значение директории для файлов кеша.
func (cfg *impl) DirectoryCache() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.CacheDirectory)
}

// DirectoryConfig Значение директории для подключаемых или дополнительных конфигураций приложения.
func (cfg *impl) DirectoryConfig() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.ConfigDirectory)
}

// DirectoryLog Значение директории для файлов журнала.
func (cfg *impl) DirectoryLog() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.LogDirectory)
}

// FileConfig Значение пути и имени конфигурационного файла приложения.
func (cfg *impl) FileConfig() string {
	return cfg.AbsolutePathAndUpdate(&cfg.bootstrapConfiguration.ConfigFile)
}

// FilePid Значение пути и имени PID файла приложения.
func (cfg *impl) FilePid() string { return cfg.bootstrapConfiguration.PidFile }

// FileState Значение пути и имени файла хранения состояния приложения.
func (cfg *impl) FileState() string { return cfg.bootstrapConfiguration.StateFile }

// FileSocket Значение пути и имени сокет файла коммуникаций с приложением.
func (cfg *impl) FileSocket() string { return cfg.bootstrapConfiguration.SocketFile }

// FileLogWithName Создание или открытие на запись файла журнала, в директории для журналов, с указанием
// названия файла журнала. Возвращается писатель, который необходимо закрыть, после окончания записи данных в журнал.
func (cfg *impl) FileLogWithName(logName string) (ret io.WriteCloser, err error) {
	const dirMode, fileMode = os.FileMode(0750), os.FileMode(0640)
	var (
		directoryName string
		fileName      string
		fileHandle    *os.File
	)

	directoryName = cfg.DirectoryLog()
	switch err = os.MkdirAll(directoryName, dirMode); {
	case err == nil:
	default:
		return
	}
	fileName = path.Join(cfg.DirectoryLog(), logName)
	if fileHandle, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode); err == nil {
		ret = fileHandle
	}

	return
}

// ConfigurationUnionSprintf Печать объединённой конфигурации приложения в строку.
func (cfg *impl) ConfigurationUnionSprintf() (ret string) {
	const structName = "*struct.UnionConfiguration"

	if cfg.conf != nil && cfg.conf.Union != nil {
		ret = rexUnionStructureHeader.ReplaceAllString(debug.DumperString(cfg.conf.Union), structName)
	}

	return
}

// ConfigurationByType Возвращает объект конфигурации соответствующий указанному типу объекта.
// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
func (cfg *impl) ConfigurationByType(t reflect.Type) (any, error) {
	return cfg.ConfigurationByTypeName(t.String())
}

// ConfigurationByTypeName Возвращает объект конфигурации соответствующий указанному названию типа объекта.
// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
func (cfg *impl) ConfigurationByTypeName(typeName string) (ret any, err error) {
	var (
		n     int
		found bool
	)

	for n = range cfg.conf.Items {
		if typeName == cfg.conf.Items[n].Type.String() {
			ret, found = cfg.conf.Items[n].Original, true
			break
		}
	}
	if !found {
		ret, err = nil, cfg.Errors().ConfigurationObjectNotFound.Bind(typeName)
		return
	}

	return
}

// ConfigurationByObject Возвращает объект конфигурации соответствующий типу переданного объекта, сам переданный
// объект никак не изменяется, он служит только для определения типа данных.
// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
func (cfg *impl) ConfigurationByObject(o any) (ret any, err error) {
	var rt reflect.Type

	if _, rt, err = reflectStructObject(o); err != nil {
		return
	}
	ret, err = cfg.ConfigurationByType(rt)

	return
}

// ConfigurationCopyByObject Если существует конфигурация с типом данных идентичным переданному объекту,
// тогда данные конфигурации копируются в переданный объект.
// Если объект конфигурации с указанным типом не регистрировался, будет возвращена ошибка.
// Объект должен передаваться по адресу, иначе его заполнение не возможно и будет возвращена ошибка.
func (cfg *impl) ConfigurationCopyByObject(o any) (err error) {
	var (
		rt  reflect.Type
		src any
	)

	// Определение типа объекта.
	if _, rt, err = reflectStructObject(o); err != nil {
		return
	}
	// Получение объекта из конфигурации.
	if src, err = cfg.ConfigurationByType(rt); err != nil {
		return
	}
	// Копирование данных их объекта конфигурации в объект результата.
	if err = kitModuleCfgCpy.All(o, src); err != nil {
		err = cfg.Errors().ConfigurationObjectCopy.Bind(rt.String(), err)
		return
	}

	return
}
