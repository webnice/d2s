package web

import (
	"crypto/tls"
	"net"
)

// ListenAndServe Открытие адреса или сокета без использования конфигурации веб сервера (конфигурация по
// умолчанию), запуск веб сервера для обслуживания входящих соединений.
func (web *impl) ListenAndServe(addr string) Interface {
	var conf *Configuration

	if conf, web.err = parseAddress(addr); web.err != nil {
		return web
	}
	return web.ListenAndServeWithConfig(conf)
}

// ListenAndServeTLS Открытие адреса или сокета с использованием TLS, без использования конфигурации веб сервера
// (конфигурация по умолчанию), запуск веб сервера для обслуживания входящих соединений.
func (web *impl) ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface {
	var conf *Configuration

	if conf, web.err = parseAddress(addr); web.err != nil {
		return web
	}
	conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM = certFile, keyFile

	return web.ListenAndServeTLSWithConfig(conf, tlsConfig)
}

// ListenAndServeWithConfig Настройка сервера с использованием переданной конфигурации, открытие адреса или сокета
// на прослушивание, запуск веб сервера для обслуживания входящих соединений.
func (web *impl) ListenAndServeWithConfig(conf *Configuration) Interface {
	var listener net.Listener

	if conf == nil {
		web.err = Errors().NoConfiguration()
		return web
	}
	web.cfg = conf
	if listener, web.err = web.NewListener(web.cfg); web.err != nil {
		return web
	}

	return web.ServeWithId(listener, web.cfg.ID)
}

// ListenAndServeTLSWithConfig Настройка сервера с использованием переданной конфигурации в режиме TLS, открытие
// адреса или сокета на прослушивание, запуск веб сервера для обслуживания входящих соединений.
func (web *impl) ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface {
	var listener net.Listener

	if conf == nil {
		web.err = Errors().NoConfiguration()
		return web
	}
	web.cfg = conf
	if tlsConfig == nil {
		tlsConfig, web.err = web.net.NewTLSConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM)
		if web.err != nil {
			return web
		}
	}
	if listener, web.err = web.NewListenerTLS(web.cfg, tlsConfig); web.err != nil {
		return web
	}

	return web.ServeTLSWithId(listener, tlsConfig, web.cfg.ID)
}

// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов переданных в процесс веб сервера из systemd.
func (web *impl) ListenersSystemdWithoutNames() ([]net.Listener, error) {
	return web.listenersSystemdWithoutNames()
}

// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов переданных в процесс веб сервера
// из systemd.
func (web *impl) ListenersSystemdWithNames() (map[string][]net.Listener, error) {
	return web.listenersSystemdWithNames()
}

// ListenersSystemdTLSWithoutNames Возвращает срез net.listener для TLS сокетов переданных в процесс веб сервера
// из systemd.
func (web *impl) ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) ([]net.Listener, error) {
	return web.listenersSystemdTLSWithoutNames(tlsConfig)
}

// ListenersSystemdTLSWithNames Возвращает карту срезов net.listener для TLS сокетов переданных в процесс веб сервера
// из systemd.
func (web *impl) ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (map[string][]net.Listener, error) {
	return web.listenersSystemdTLSWithNames(tlsConfig)
}

// NewListener Создание нового слушателя соединений net.Listener на основе конфигурации веб сервера.
func (web *impl) NewListener(conf *Configuration) (ret net.Listener, err error) {
	if web.net.IsRunning() {
		err = Errors().AlreadyRunning()
		return
	}
	ret, _, err = web.net.NewListener(&conf.Configuration)

	return
}

// NewListenerTLS Создание нового слушателя соединений net.Listener в режиме TLS, на основе конфигурации
// веб сервера.
func (web *impl) NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, err error) {
	if web.net.IsRunning() {
		err = Errors().AlreadyRunning()
		return
	}
	ret, _, err = web.net.NewListenerTLS(&conf.Configuration, tlsConfig)

	return
}
