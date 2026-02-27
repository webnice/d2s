package server

import "time"

// GrpcConfiguration Структура конфигурации GRPC сервера.
type GrpcConfiguration struct {
	// Server Конфигурация WEB сервера.
	Server GrpcServerConfiguration `yaml:"Server"`
}

// GrpcServers Структура конфигурации группы GRPC серверов.
type GrpcServers struct {
	GrpcServers []GrpcConfiguration `yaml:"GRPCServers"`
}

// GrpcServerConfiguration Структура конфигурации веб сервера.
type GrpcServerConfiguration struct {
	// TLSPublicKeyPEM Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA сертификаты
	// всех промежуточных центров сертификации, если ими подписан ключ.
	TLSPublicKeyPEM string `yaml:"TLSPublicKeyPEM" json:"tls_public_key_pem"`

	// TLSPrivateKeyPEM Путь и имя файла содержащего приватный ключ в PEM формате.
	TLSPrivateKeyPEM string `yaml:"TLSPrivateKeyPEM" json:"tls_private_key_pem"`

	// Host IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов.
	// Default value: "0.0.0.0"
	Host string `yaml:"Host" json:"host" default-value:"0.0.0.0"`

	// Port tcp/ip порт занимаемый сервером.
	// Default value: 80
	Port uint16 `yaml:"Port" json:"port" default-value:"80"`

	// Socket Unix socket на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac.
	// Default value: "" - unix socket is off
	Socket string `yaml:"Socket" json:"socket" default-value:"-"`

	// Mode Режим работы, tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
	// tcp - Сервер поднимается на указанном Host:Port.
	// socket - Сервер поднимается на socket, только для unix-like операционных систем.
	// systemd - Порт или сокет открывает systemd и передаёт слушателя порта через файловый дескриптор сервису
	//           запущенному от пользователя без права открытия привилегированных портов.
	// Default value: "tcp"
	Mode string `yaml:"Mode" json:"mode" default-value:"tcp"`

	// WriteBufferSize Размер буфера записи, определяет максимальное количество данных для одного системного вызова.
	// Памяти выделяется в 2 раза больше размера буфера записи.
	// Если указано значение 0 - это отключит буфер записи.
	// Default value: 32768 (32 килобайта).
	WriteBufferSize int `yaml:"WriteBufferSize"`

	// ReadBufferSize Размер буфера чтения, определяет максимальное читаемых количество данных для одного системного
	// вызова.
	// Если указано значение 0 - это отключит буфер чтения.
	// Default value: 32768 (32 килобайта).
	ReadBufferSize int `yaml:"ReadBufferSize"`

	// InitialWindowSize Размер окна для потоковых данных.
	// Минимальное значение 65536 байт или 64 килобайта, любое меньшее значение будет проигнорировано.
	// Default value: 65536 (64 килобайта).
	InitialWindowSize int32 `yaml:"InitialWindowSize"`

	// InitialConnWindowSize Размер окна для данных запросов (подключения).
	// Минимальное значение 65536 байт или 64 килобайта, любое меньшее значение будет проигнорировано.
	// Default value: 65536 (64 килобайта).
	InitialConnWindowSize int32 `yaml:"InitialConnWindowSize"`

	// KeepaliveMaxConnectionIdle Максимальное время простоя не занятого соединения, по истечении которого
	// соединение закрывается путём отправки сообщения об отказе.
	// Продолжительность простоя соединения определяется с момента установки соединения либо с момента времени,
	// когда количество не занятых соединений стало равным нулю.
	// Default value: 0s - бесконечность.
	KeepaliveMaxConnectionIdle time.Duration `yaml:"KeepaliveMaxConnectionIdle"`

	// KeepaliveMaxConnectionAge Максимальное значение времени в течении которого соединени может оставаться открытым.
	// Значение по умолчанию - бесконечность.
	// Для значения боьше 0, добавляется случайное значение в размере +/- 10%.
	// Default value: 0s - бесконечность.
	KeepaliveMaxConnectionAge time.Duration `yaml:"KeepaliveMaxConnectionAge"`

	// KeepaliveMaxConnectionAgeGrace Дополнительная отсрочка принудительного закрытия соединения после истечения
	// значения определенного в MaxConnectionAge.
	// Значение по умолчанию - бесконечность.
	// Default value: 0s - бесконечность.
	KeepaliveMaxConnectionAgeGrace time.Duration `yaml:"KeepaliveMaxConnectionAgeGrace"`

	// KeepaliveTime Время, по истечении которого, если сервер не видит никакой активности, выполняется
	// проверка соединения с клиентом.
	// Минимальное значение 1 секунда, если указано значение меньше 1 секунды, используется 1 секунда.
	// Значение по умолчанию - 2 часа.
	// Default value: 2h
	KeepaliveTime time.Duration `yaml:"KeepaliveTime"`

	// KeepaliveTimeout Время ожидания ответа на сообщение пинга, по истечении которого, соединение закрывается.
	// Значение по умолчанию 20 секунд.
	// Default value: 20s
	KeepaliveTimeout time.Duration `yaml:"KeepaliveTimeout"`

	// KeepaliveMinTime Минимальное количество времени которое клиент должен подождать перед отправкой пинг запроса.
	// Default value: 5m
	KeepaliveMinTime time.Duration `yaml:"KeepaliveMinTime"`

	// KeepalivePermitWithoutStream Разрешение поддержки пинг-запросов при отсутствии активных потоков.
	// Если указано "ложь", при отсутствии активных потоков, на запрос пинг от клиента, сервер ответит
	// сообщением GO AWAY и закроет соединение.
	// Default value: false
	KeepalivePermitWithoutStream bool `yaml:"KeepalivePermitWithoutStream"`

	// MaxRecvMsgSize Максимальный размер сообщения в байтах, которое может принимать сервер.
	// Если не указано, значение по умолчанию 4 мегабайта.
	// Default value: 4194304 (4 мегабайта).
	MaxRecvMsgSize int `yaml:"MaxRecvMsgSize"`

	// MaxSendMsgSize Максимальный размер сообщения в байтах, которое может отправить сервер.
	// Если не указано, значение по умолчанию равно максимальному числу int32 = 2147483647 или 2 гигабайта.
	// Default value: 2147483647 (=math.MaxInt32)
	MaxSendMsgSize int `yaml:"MaxSendMsgSize"`

	// MaxConcurrentStreams Максимальное количество потоковых соединений, которое может принимать сервер.
	// Default value: 0 - нет ограничений.
	MaxConcurrentStreams uint32 `yaml:"MaxConcurrentStreams"`

	// ConnectionTimeout Время ожидания установки соединения.
	// Нулевое или отрицательное значение приведёт к немедленному завершению соединения, которое не
	// успеет установиться и будет приводить к ошибке.
	// Default value: 120s
	ConnectionTimeout time.Duration `yaml:"ConnectionTimeout"`

	// MaxHeaderListSize Максимальный, не сжатый, размер заголовков, принимаемый сервером.
	// Default value: 0 - без ограничений.
	MaxHeaderListSize uint32 `yaml:"MaxHeaderListSize"`

	// HeaderTableSize Размер динамической таблицы заголовков для потоковой передачи данных.
	// Default value: 0 - без ограничений.
	HeaderTableSize uint32 `yaml:"HeaderTableSize"`

	// NumStreamWorkers Количество потоков, которые держит сервер, для обработки входящих запросов.
	// Значение равное нулю заставляет сервер запускать новый поток для каждого нового запроса.
	// Default value: 0 - нет потоков, для каждого нового запроса создаётся новый поток.
	NumStreamWorkers uint32 `yaml:"NumStreamWorkers"`
}
