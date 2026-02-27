package web

import (
	"time"

	"github.com/webnice/net"
)

// Configuration Структура конфигурации веб сервера.
type Configuration struct {
	// Описание сетевого или локального доступа к серверу.
	net.Configuration `yaml:"Network" json:"network"`

	// TODO Сделать ограничение по доменам.
	// Domains Список всех доменов, на которые отвечает сервер.
	// Если не пусто, то для всех других доменов будет ответ "Requested host unavailable".
	// Default value: [] - all domain
	Domain []string `yaml:"Domain" json:"domain"`

	// ReadTimeout Время ожидания запроса включая ReadHeaderTimeout.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	ReadTimeout time.Duration `yaml:"ReadTimeout" json:"read_timeout"`

	// ReadHeaderTimeout Время ожидания заголовка запроса.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	ReadHeaderTimeout time.Duration `yaml:"ReadHeaderTimeout" json:"read_header_timeout"`

	// WriteTimeout Время ожидания выдачи ответа.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	WriteTimeout time.Duration `yaml:"WriteTimeout" json:"write_timeout"`

	// IdleTimeout Максимальное время ожидания следующего входящего соединения для открытого сокета, до его закрытия.
	// Используется при включённом keep-alives.
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	IdleTimeout time.Duration `yaml:"IdleTimeout" json:"idle_timeout"`

	// ShutdownTimeout Максимальное время ожидания завершения работы веб сервера до начала принудительного обрыва
	// соединений и остановки процессов.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 30s
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout" json:"shutdown_timeout" default-value:"30s"`

	// MaxHeaderBytes Максимальный размер заголовка запроса.
	// Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
	MaxHeaderBytes int `yaml:"MaxHeaderBytes" json:"max_header_bytes" default-value:"1048576"`

	// TODO Сделать ограничение на максимальный размер тела запроса.
	// MaxBodyBytes Максимальный размер тела запроса.
	// Default value: 0 - unlimited
	MaxBodyBytes uint64 `yaml:"MaxBodyBytes" json:"max_body_bytes"`

	// KeepAliveDisable Отключение режима "оставаться в живых" (keep alive).
	// Default value: false - keep alive are enabled.
	KeepAliveDisable bool `yaml:"KeepAliveDisable" json:"keep_alive_disable"`

	// DisableGeneralOptionsHandler Отключение обработки HTTP запросов методом OPTION.
	// Если установлено значение "истина", на все запросы методом OPTION сервер отвечает кодом 200 и передаёт в
	// заголовке Content-Length значение 0.
	// Если установлено значение "ложь", запросы передаются контроллеру и обрабатываются обычным образом.
	// Default value: false - запросы передаются контроллеру.
	DisableGeneralOptionsHandler bool `yaml:"DisableGeneralOptionsHandler" json:"disable_general_options_handler"`
}

/**

   Пример конфигурации YAML:

      ## Описание сетевого или локального доступа к серверу.
      Network:
            ## ID Уникальный идентификатор сервера, может быть любым уникальным строковым значением.
            ## Если при запуске сервера значение не указано, создаётся уникальное временное значение, меняющееся
            ## при каждом запуске.
            ## Default value: ""
            ID string `yaml:"ID" json:"id"`

            ## Публичный адрес по которому сервер доступен извне.
            ## Например, если сервер находится за прокси, тут указывается реальный адрес подключения к серверу.
            ## Default value: ""
            Address: !!str "http://localhost/"

            ## IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов.
            ## Default value: "0.0.0.0".
            #Host: !!str "[2a03:e2c0:a32::2]"
            #Host: !!str "example.hostname.local"
            Host: !!str "0.0.0.0"

            ## Tcp/ip порт занимаемый сервером.
            ## Default value: 80
            Port: !!int 80

            ## Юникс сокет на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac.
            ## Default value: ""
            Socket: !!str "run/example.sock"

            ## Файловые разрешения доступа к юникс-сокету.
            ## Значение задаётся в восьмеричной системе и не должно превышать 32 бита.
            ## Default value: "0666"
            SocketMode: !!str "0666"

            ## Режим открытия сокета, возможные значения: tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
            ## udp, udp4, udp6 - Сервер поднимается на указанном Host:Port;
            ## tcp, tcp4, tcp6 - Сервер поднимается на указанном Host:Port;
            ## unix, unixpacket - Сервер поднимается на указанном unix/unixpacket;
            ## socket  - Сервер поднимается на socket, только для unix-like операционных систем. Параметры Host:Port
            ##           игнорируются, используется только путь к сокету;
            ## systemd - Порт или сокет открывает systemd и передаёт слушателя порта через файловый дескриптор сервису,
            ##           запущенному от пользователя без права открытия привилегированных портов. Максимально удобный
            ##           способ при использовании правильного безопасно настроенного linux сервера.
            ##           Более подробно можно посмотреть в документации man systemd.socket(5);
            ## Default value: "tcp"
            Mode: !!str "tcp"

            ## Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA
            ## сертификаты всех промежуточных центров сертификации, если ими подписан ключ.
            ## Применяется только для TCP соединений, для UDP не используется.
            ## Default value: ""
            TLSPublicKeyPEM: !!str "/etc/application/certificate.pub"

            ## Путь и имя файла содержащего секретный/приватный ключ в PEM формате.
            ## Применяется только для TCP соединений, для UDP не используется.
            ## Default value: ""
            TLSPrivateKeyPEM: !!str "/etc/application/certificate.key"

            ## Включение прокси-протокола.
            ## Прокси-протокол позволяет серверу получать информацию о подключении клиента, передаваемую через
            ## прокси-серверы и средства балансировки нагрузки, такие как Nginx, HAProxy, Amazon Elastic Load
            ## Balancer (ELB) и многие другие.
            ## Поддерживаются запросы с реализацией прокси протокола версий 1 и 2.
            ## PROXY protocol: https://www.haproxy.org/download/2.3/doc/proxy-protocol.txt.
            ## Default value: false
            ProxyProtocol: !!bool false

            ## Максимальное время ожидания получения данных о клиенте через прокси-протокол.
            ## Время ожидание используется только при включённом прокси-протоколе.
            ## Default value: 0s - no timeout
            ProxyProtocolReadHeaderTimeout: 0s

      ## Список всех доменов, на которые отвечает сервер.
      ## Если не пусто, то для всех других доменов будет ответ "Requested host unavailable".
      ## Default value: [] - любое имя домена.
      Domain:
      - !!str "localhost"
      - !!str "example.domain.tld"

      ## Время ожидания запроса включая ReadHeaderTimeout.
      ## Если не указано или рано 0 - таймаута нет.
      ## Default value: 0 - no timeout
      ReadTimeout: 0s

      ## Время ожидания заголовка запроса.
      ## Если не указано или рано 0 - таймаута нет.
      ## Default value: 0 - no timeout
      ReadHeaderTimeout: 0s

      ## Время ожидания выдачи ответа.
      ## Если не указано или рано 0 - таймаута нет.
      ## Default value: 0 - no timeout
      WriteTimeout: 0s

      ## Максимальное время ожидания следующего входящего соединения для открытого сокета, до его закрытия.
      ## Используется при включённом keep-alives.
      ## Если не указано или рано 0 - таймаута нет
      ## Default value: 0 - no timeout
      IdleTimeout: 0s

      ## Максимальное время ожидания завершения работы веб сервера до начала принудительного обрыва
      ## соединений и остановки процессов.
      ## Если не указано или рано 0 - таймаута нет.
      ## Default value: 30s
      ShutdownTimeout: 30s

      ## Максимальный размер заголовка запроса.
      ## Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
      MaxHeaderBytes: !!int 1048576

      ## Максимальный размер тела запроса.
      ## Default value: 0 - без ограничений.
      MaxBodyBytes: !!int 0

      ## Отключение режима "оставаться в живых" (keep alive).
      ## Default value: false - оставаться в живых включён.
      KeepAliveDisable: !!bool false

      ## Отключение обработки HTTP запросов методом OPTION.
      ## Если установлено значение "истина", на все запросы методом OPTION сервер отвечает кодом 200 и передаёт в
      ## заголовке Content-Length значение 0.
      ## Если установлено значение "ложь", запросы передаются контроллеру и обрабатываются обычным образом.
      ## Default value: false - запросы передаются контроллеру.
      DisableGeneralOptionsHandler: !!bool false


**/
