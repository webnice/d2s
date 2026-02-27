package web

import (
	"crypto/tls"
	"net"
	"net/http"

	wnet "github.com/webnice/net"
)

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New() Interface {
	var web = &impl{
		net: wnet.New(),
	}
	web.listenersSystemdWithoutNames = web.net.ListenersSystemdWithoutNames
	web.listenersSystemdWithNames = web.net.ListenersSystemdWithNames
	web.listenersSystemdTLSWithoutNames = web.net.ListenersSystemdTLSWithoutNames
	web.listenersSystemdTLSWithNames = web.net.ListenersSystemdTLSWithNames

	return web
}

// ID Уникальный идентификатор сервера.
// Если идентификатор не был указан в конфигурации, создаётся при запуске.
func (web *impl) ID() string { return web.net.ID() }

// Handler Назначение обработчика запросов ВЕБ сервера.
// Обработчик необходимо назначить до запуска ВЕБ сервера.
func (web *impl) Handler(handler http.Handler) Interface { web.handler = handler; return web }

// Clean Очистка последней ошибки.
func (web *impl) Clean() Interface { web.err = nil; web.net.Clean(); return web }

// Errors Справочник ошибок.
func (web *impl) Errors() *Error { return errSingleton }

// Error Функция возвращает последнюю ошибку веб сервера или библиотеки "github.com/webnice/net", на которой
// основан ВЕБ сервер.
func (web *impl) Error() error {
	switch web.err {
	case nil:
		return web.net.Error()
	default:
		return web.err
	}
}

// Serve Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener.
func (web *impl) Serve(ltn net.Listener) Interface { return web.ServeTLS(ltn, nil) }

// ServeWithId Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// указанием ID сервера.
func (web *impl) ServeWithId(ltn net.Listener, id string) Interface {
	return web.ServeTLSWithId(ltn, nil, id)
}

// ServeTLS Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// использованием TLS.
func (web *impl) ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface {
	return web.ServeTLSWithId(ltn, tlsConfig, "")
}

// ServeTLSWithId Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// использованием TLS и указанием ID сервера.
func (web *impl) ServeTLSWithId(ltn net.Listener, tlsConfig *tls.Config, id string) Interface {
	var (
		listener net.Listener
		isTls    bool
	)

	if web.cfg == nil {
		web.cfg, web.err = parseAddress(ltn.Addr().String())
		defaultConfiguration(web.cfg)
	}
	web.cfg.ID, web.cfg.Configuration.ID = id, id
	if web.server, isTls = web.makeServer(tlsConfig); web.Error() != nil {
		return web
	}
	if listener = ltn; isTls {
		listener = tls.NewListener(ltn, tlsConfig)
	}
	web.net.Handler(web.server.Serve)
	web.net.ServeWithId(listener, web.cfg.ID)

	return web
}

// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
// Если сервер не запущен, функция завершается немедленно.
func (web *impl) Wait() Interface { web.net.Wait(); return web }

// Stop Отправка сигнала прерывания работы веб сервера с учётом значения ShutdownTimeout.
func (web *impl) Stop() Interface { web.net.Stop(); return web }
