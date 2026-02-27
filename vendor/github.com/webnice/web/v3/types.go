package web

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/webnice/dic"
	wnet "github.com/webnice/net"
)

const (
	netTcp        = "tcp"
	netTcp4       = "tcp4"
	netTcp6       = "tcp6"
	netUnix       = "unix"
	netUnixPacket = "unixpacket"
	netSocket     = "socket"
	netSystemd    = "systemd"
)

// Справочники реализованы в отдельной библиотеке. Но есть две причины их появления тут:
// 1. При написании кода не требуется подключать отдельную библиотеку, а библиотека web уже будет подключена.
// 2. Все справочник уже являются объектами-одиночками, поэтому ссылки на них не влияют на память, только на удобство.
var (
	// Mime Справочник MIME типов.
	Mime = dic.Mime()

	// Method Справочник HTTP методов запросов.
	Method = dic.Method()

	// Header Справочник заголовков.
	Header = dic.Header()

	// Status Справочник статусов HTTP ответов.
	Status = dic.Status()
)

var _, _, _, _ = Mime, Method, Header, Status

// Объект сущности, реализующий интерфейс Interface.
type impl struct {
	err     error          // Сохранение последней ошибки.
	net     wnet.Interface // Интерфейс "github.com/webnice/net".
	cfg     *Configuration // Конфигурация веб сервера.
	handler http.Handler   // Обработчик запросов.
	server  *http.Server   // Объект веб сервера "net/http".

	// Функции которые невозможно протестировать в обычных условиях практически никаким образом.
	listenersSystemdWithoutNames    func() ([]net.Listener, error)
	listenersSystemdWithNames       func() (map[string][]net.Listener, error)
	listenersSystemdTLSWithoutNames func(*tls.Config) ([]net.Listener, error)
	listenersSystemdTLSWithNames    func(*tls.Config) (map[string][]net.Listener, error)
}
