package types

import kmll "github.com/webnice/kit/v4/module/log/level"

// BootstrapConfiguration Структура стартовой конфигурации приложения.
// Параметры тегов "env", "default", "help", описанные в блоке kong, обрабатываются в момент инициализации
// CLI (первыми), и имеют наивысший приоритет.
// Параметры тегов "env-name", "default-value", "description" предназначены для конфигурационного файла и
// обрабатываются после загрузки конфигурационного файла (вторыми).
// Если значение переменной, после инициализации CLI, равно пустому значению или значению по умолчанию, тогда, после
// загрузки конфигурационного файла, значение устанавливается согласно приоритету:
// 1. Значение указанное в конфигурационном файле (наивысший приоритет).
// 2. Значение указанное в переменной окружения с именем указанным в "env-name".
// 3. Значения полученные через интерфейс types.ConfigurationDefaulter, если структура реализует этот интерфейс.
// 4. Значение указанное в "default-value".
// 5. Значение по умолчанию, для типа golang (пустое значение).
// Если необходимо, при загрузке конфигурационного файла, загрузить значение из переменной окружения, тогда, переменная
// в конфигурационном файле, должна отсутствовать, либо должна быть определена пустым значением.
// Помимо этого, переменная может быть загружена из двух разных переменных окружения, имеющих разный приоритет
// загрузки, первая описывается в блоке kong, имя "env", обрабатывается на первом этапе, вторая описывается
// за пределами блока kong, имеет имя тега "env-name" и обрабатывается после загрузки конфигурационного файла.
type BootstrapConfiguration struct {
	ApplicationForkWorkerMaster    string     `kong:",,,,,,,,,,,,,,,,,,,,,,,,,,hidden,,env='4FD0501A-84A5-4368-B23E-9DCAFE534EDE',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,help='Параметры связи подчинённого процесса с управляющим процессом.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"-"                       env-name:"-"  default-value:"-"  description:"-"`
	ApplicationForkWorkerComponent string     `kong:",,,,,,,,,,,,,,,,,,,,,,,,,,hidden,,env='1D12FC38-F7A3-4C8F-BEC9-448BE6BDFB9F',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,help='Параметры запуска компоненты на стороне подчинённого процесса.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"-"                       env-name:"-"  default-value:"-"  description:"-"`
	ApplicationForkWorkerTarget    string     `kong:",,,,,,,,,,,,,,,,,,,,,,,,,,hidden,,env='3DEDF187-1FE5-4568-AD5D-68FD4984A5F6',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,help='Параметры запуска процесса на стороне подчинённого процесса.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"-"                       env-name:"-"  default-value:"-"  description:"-"`
	ApplicationTargetlevel         uint16     `kong:"name='target-level',,,,,,,hidden,,env='3D131B37-0ECE-46BD-861D-91DF3F739FAA',,default='${ApplicationTargetlevel}',,help='Целевой уровень выполнения приложения, по умолчанию 65535.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"ApplicationTargetlevel"  env-name:"-"  default-value:"-"  description:"Целевой уровень выполнения приложения.\nПо умолчанию 65535."`
	ApplicationDebug               bool       `kong:"name='debug',,,,,,,,,,,,,,,,,,,,,,env='FC53AACC-ED71-4446-A95E-655BC03447B6',,default='${ApplicationDebug}',,,,,,,,help='Включение режима отладки приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"ApplicationDebug"        env-name:"-"  default-value:"-"  description:"Включение режима отладки приложения."`
	ApplicationName                string     `kong:"name='application-name',,,hidden,,env='C5351023-07D6-4FCD-9F07-E97C94D7F697',,default='${ApplicationName}',,,,,,,,,help='Название приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"ApplicationName"         env-name:"-"  default-value:"-"  description:"Название приложения."`
	HomeDirectory                  string     `kong:"name='home-directory',,,,,hidden,,env='212C3FE9-1CD3-4DCD-989E-712319FE0CEF',,default='${HomeDirectory}',,,,,,,,,,,help='Домашняя директория приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"HomeDirectory"           env-name:"-"  default-value:"-"  description:"Домашняя директория приложения."`
	WorkingDirectory               string     `kong:"name='working-directory',,,,,,,,,,env='97723E20-0FC3-42BD-A210-941D47A23CAE',,default='${WorkingDirectory}',,,,,,,,help='Рабочая директория приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"WorkingDirectory"        env-name:"-"  default-value:"-"  description:"Рабочая директория приложения."`
	TempDirectory                  string     `kong:"name='temp-directory',,,,,hidden,,env='C9BBD482-508C-4857-94D9-1CC3CAE38E63',,default='${TempDirectory}',,,,,,,,,,,help='Директория для временных файлов.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"TempDirectory"           env-name:"-"  default-value:"-"  description:"Директория для временных файлов."`
	CacheDirectory                 string     `kong:"name='cache-directory',,,,hidden,,env='F055F1EA-5CF0-4D50-956B-6A1FF200B076',,default='${CacheDirectory}',,,,,,,,,,help='Директория для файлов кеша.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"CacheDirectory"          env-name:"-"  default-value:"-"  description:"Директория для файлов кеша."`
	ConfigDirectory                string     `kong:"name='config-directory',,,hidden,,env='CC575694-2AAC-488A-928F-D54DEBA22ED8',,default='${ConfigDirectory}',,,,,,,,,help='Директория для подключаемых или дополнительных конфигураций приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"ConfigDirectory"         env-name:"-"  default-value:"-"  description:"Директория для подключаемых или дополнительных конфигураций приложения."`
	LogDirectory                   string     `kong:"name='log-directory',,,,,,hidden,,env='CAF4DC58-C6DE-46AB-9FE2-BA5A44BA63AF',,default='${LogDirectory}',,,,,,,,,,,,help='Директория для файлов журнала приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"LogDirectory"            env-name:"-"  default-value:"-"  description:"Директория для файлов журнала приложения."`
	ConfigFile                     string     `kong:"name='conf',,,,,,,,,,,,,,,,,,,,,,,env='BE735413-0146-4AE5-945C-F378719D8A2D',,default='${ConfigFile}',,,,,,,,,,,,,,help='Путь и имя файла конфигурации приложения, если не указан, приложение ищет конфигурацию самостоятельно.',," yaml:"-"                       env-name:"-"  default-value:"-"  description:"-"`
	PidFile                        string     `kong:"name='pid',,,,,,,,,,,,,,,,,,,,,,,,env='C9E49D1C-D1F2-4C16-BCE3-0412890F8443',,default='${PidFile}',,,,,,,,,,,,,,,,,help='Путь и имя PID файла приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"PidFile"                 env-name:"-"  default-value:"-"  description:"Путь и имя PID файла приложения."`
	StateFile                      string     `kong:"name='state',,,,,,,,,,,,,,,,,,,,,,env='5AC281D5-349C-4EE7-AA91-C897043B3EB5',,default='${StateFile}',,,,,,,,,,,,,,,help='Путь и имя файла хранения состояния приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"StateFile"               env-name:"-"  default-value:"-"  description:"Путь и имя файла хранения состояния приложения."`
	SocketFile                     string     `kong:"name='socket',,,,,,,,,,,,,,,,,,,,,env='D5246FFD-748A-4A4C-AF66-2B840281528F',,default='${SocketFile}',,,,,,,,,,,,,,help='Сокет файл коммуникаций с приложением, только для *nix систем, путь и имя файла.',,,,,,,,,,,,,,,,,,,,,,,," yaml:"SocketFile"              env-name:"-"  default-value:"-"  description:"Сокет файл коммуникаций с приложением, только для *nix систем, путь и имя файла."`
	LogLevel                       kmll.Level `kong:"name='log-level',,,,,,,,,,hidden,,env='D949FEDA-B6C5-4FB6-A2C7-ABFEAA99473B',,default='${LogLevel}',,,,,,,,,,,,,,,,help='Уровень логирования по умолчанию до загрузки конфигурации приложения.',,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,," yaml:"LogLevel"                env-name:"-"  default-value:"-"  description:"Уровень логирования по умолчанию до загрузки конфигурации приложения."`
}

// CommandGroup Описание группы команд отображаемое при выводе помощи к командам.
type CommandGroup struct {
	// Key Идентификатор группы команд, ключ связи описания группы с командой, входящей в группу.
	Key string

	// Title Заголовок группы команд, может быть пустым.
	Title string

	// Description Описание группы команд. Не может быть пустым, если пустое, тогда описание группы игнорируется.
	Description string
}

// BootstrapDefaultValueFn Функция, возвращающая значение по умолчанию для параметра конфигурации приложения.
type BootstrapDefaultValueFn func() string

// BootstrapCommandGroupValueFn Функция возвращает список групп команд, отображаемый при выводе помощи.
type BootstrapCommandGroupValueFn func() []CommandGroup

// BootstrapDefaultValue Описание функций которые вернут значения по умолчанию для конфигурации приложения.
type BootstrapDefaultValue struct {
	ApplicationTargetlevel BootstrapDefaultValueFn      // Целевой уровень выполнения приложения.
	ApplicationDebug       BootstrapDefaultValueFn      // Режим отладки приложения.
	ApplicationName        BootstrapDefaultValueFn      // Название приложения.
	ApplicationDescription BootstrapDefaultValueFn      // Описание приложения, отображаемое в строке помощи.
	HomeDirectory          BootstrapDefaultValueFn      // Домашняя директория приложения.
	WorkingDirectory       BootstrapDefaultValueFn      // Рабочая директория приложения.
	TempDirectory          BootstrapDefaultValueFn      // Директория для временных файлов.
	CacheDirectory         BootstrapDefaultValueFn      // Директория для файлов кеша.
	ConfigDirectory        BootstrapDefaultValueFn      // Директория для подключаемых или дополнительных конфигураций приложения.
	LogDirectory           BootstrapDefaultValueFn      // Директория для файлов журнала приложения.
	ConfigFile             BootstrapDefaultValueFn      // Путь и имя файла конфигурации приложения.
	PidFile                BootstrapDefaultValueFn      // Путь и имя PID файла приложения.
	StateFile              BootstrapDefaultValueFn      // Путь и имя файла хранения состояния приложения.
	SocketFile             BootstrapDefaultValueFn      // Сокет файл коммуникаций с приложением.
	LogLevel               BootstrapDefaultValueFn      // Уровень логирования по умолчанию до загрузки конфигурации приложения.
	CommandGroup           BootstrapCommandGroupValueFn // Список групп команд.
}

// BootstrapConfigurationForkWorker Структура параметров для работы в режиме forkWorker.
type BootstrapConfigurationForkWorker struct {
	// Master Параметры связи подчинённого процесса с управляющим процессом.
	Master string

	// Component Параметры запуска компоненты на стороне подчинённого процесса.
	Component string

	// Target Параметры запуска процесса на стороне подчинённого процесса.
	Target string
}
