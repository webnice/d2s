package _const

const (
	// EnvironmentApplicationForkWorkerMaster Константа для поиска переменной окружения отвечающей за:
	// Параметры связи подчинённого процесса с управляющим процессом.
	EnvironmentApplicationForkWorkerMaster = "4FD0501A-84A5-4368-B23E-9DCAFE534EDE"

	// EnvironmentApplicationForkWorkerComponent Константа для поиска переменной окружения отвечающей за:
	// Параметры запуска компоненты на стороне подчинённого процесса.
	EnvironmentApplicationForkWorkerComponent = "1D12FC38-F7A3-4C8F-BEC9-448BE6BDFB9F"

	// EnvironmentApplicationForkWorkerTarget Константа для поиска переменной окружения отвечающей за:
	// Параметры запуска процесса на стороне подчинённого процесса.
	EnvironmentApplicationForkWorkerTarget = "3DEDF187-1FE5-4568-AD5D-68FD4984A5F6"

	// EnvironmentApplicationTargetlevel Константа для поиска переменной окружения отвечающей за:
	// Целевой уровень выполнения приложения, по умолчанию 65535.
	EnvironmentApplicationTargetlevel = "3D131B37-0ECE-46BD-861D-91DF3F739FAA"

	// EnvironmentApplicationDebug Константа для поиска переменной окружения отвечающей за:
	// Включение режима отладки приложения.
	EnvironmentApplicationDebug = "FC53AACC-ED71-4446-A95E-655BC03447B6"

	// EnvironmentApplicationName Константа для поиска переменной окружения отвечающей за:
	// Название приложения.
	EnvironmentApplicationName = "C5351023-07D6-4FCD-9F07-E97C94D7F697"

	// EnvironmentHomeDirectory Константа для поиска переменной окружения отвечающей за:
	// Домашняя директория приложения.
	EnvironmentHomeDirectory = "212C3FE9-1CD3-4DCD-989E-712319FE0CEF"

	// EnvironmentWorkingDirectory Константа для поиска переменной окружения отвечающей за:
	// Рабочая директория приложения.
	EnvironmentWorkingDirectory = "97723E20-0FC3-42BD-A210-941D47A23CAE"

	// EnvironmentTempDirectory Константа для поиска переменной окружения отвечающей за:
	// Директория для временных файлов.
	EnvironmentTempDirectory = "C9BBD482-508C-4857-94D9-1CC3CAE38E63"

	// EnvironmentCacheDirectory Константа для поиска переменной окружения отвечающей за:
	// Директория для файлов кеша.
	EnvironmentCacheDirectory = "F055F1EA-5CF0-4D50-956B-6A1FF200B076"

	// EnvironmentConfigDirectory Константа для поиска переменной окружения отвечающей за:
	// Директория для подключаемых или дополнительных конфигураций приложения.
	EnvironmentConfigDirectory = "CC575694-2AAC-488A-928F-D54DEBA22ED8"

	// EnvironmentLogDirectory Константа для поиска переменной окружения отвечающей за:
	// Директория для файлов журнала приложения.
	EnvironmentLogDirectory = "CAF4DC58-C6DE-46AB-9FE2-BA5A44BA63AF"

	// EnvironmentConfigFile Константа для поиска переменной окружения отвечающей за:
	// Путь и имя файла конфигурации приложения, если не указан, приложение ищет конфигурацию самостоятельно.
	EnvironmentConfigFile = "BE735413-0146-4AE5-945C-F378719D8A2D"

	// EnvironmentPidFile Константа для поиска переменной окружения отвечающей за:
	// Путь и имя PID файла приложения.
	EnvironmentPidFile = "C9E49D1C-D1F2-4C16-BCE3-0412890F8443"

	// EnvironmentStateFile Константа для поиска переменной окружения отвечающей за:
	// Путь и имя файла хранения состояния приложения.
	EnvironmentStateFile = "5AC281D5-349C-4EE7-AA91-C897043B3EB5"

	// EnvironmentSocketFile Константа для поиска переменной окружения отвечающей за:
	// Сокет файл коммуникаций с приложением, только для *nix систем, путь и имя файла.
	EnvironmentSocketFile = "D5246FFD-748A-4A4C-AF66-2B840281528F"

	// EnvironmentLogLevel Константа для поиска переменной окружения отвечающей за:
	// Уровень логирования по умолчанию до загрузки конфигурации приложения.
	EnvironmentLogLevel = "D949FEDA-B6C5-4FB6-A2C7-ABFEAA99473B"
)
