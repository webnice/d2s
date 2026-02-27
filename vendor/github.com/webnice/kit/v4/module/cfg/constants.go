package cfg

import "math"

const (
	defaultRunlevel     = uint16(10)             // Уровень выполнения компонентов приложения по умолчанию.
	defaultRunlevelExit = uint16(math.MaxUint16) // Уровень завершения работы приложения.
)

const (
	tagYaml         = "yaml"
	tagDefaultValue = "default-value"
	tagEnvName      = "env-name"
	tagDescription  = "description"
)

const (
	runeTilda    = '~'
	runeSlash    = '/'
	keySlash     = string(runeSlash)
	keyOsWindows = "windows"
)

const (
	envNameAnchorForApplicationForkWorkerMaster    = "APPLICATION_FORK_WORKER_MASTER"
	envNameAnchorForApplicationForkWorkerComponent = "APPLICATION_FORK_WORKER_COMPONENT"
	envNameAnchorForApplicationForkWorkerTarget    = "APPLICATION_FORK_WORKER_TARGET"
	envNameAnchorForApplicationTargetlevel         = "APPLICATION_TARGET_LEVEL"
	envNameAnchorForApplicationDebug               = "APPLICATION_DEBUG"
	envNameAnchorForApplicationName                = "APPLICATION_NAME"
	envNameAnchorForHomeDirectory                  = "HOME"
	envNameAnchorForWorkingDirectory               = "APPLICATION_WORKING_DIRECTORY"
	envNameAnchorForTempDirectory                  = "APPLICATION_TEMP_DIRECTORY"
	envNameAnchorForCacheDirectory                 = "APPLICATION_CACHE_DIRECTORY"
	envNameAnchorForConfigDirectory                = "APPLICATION_CONF_DIRECTORY"
	envNameAnchorForLogDirectory                   = "APPLICATION_LOG_DIRECTORY"
	envNameAnchorForConfigFile                     = "APPLICATION_CONFIGURATION"
	envNameAnchorForPidFile                        = "APPLICATION_PID"
	envNameAnchorForStateFile                      = "APPLICATION_STATE"
	envNameAnchorForSocketFile                     = "APPLICATION_SOCKET"
	envNameAnchorForLogLevel                       = "APPLICATION_LOG_LEVEL"
)

const (
	tplRunlevelExit = "Переключение значения уровня работы приложения на уровень завершения работы."
)
