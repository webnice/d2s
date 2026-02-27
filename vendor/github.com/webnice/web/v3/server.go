package web

import (
	"crypto/tls"
	"net/http"
)

// Применение настроек из конфигурации к стандартному веб серверу.
func (web *impl) makeServer(tlsConfig *tls.Config) (ret *http.Server, isTLS bool) {
	const http11 = "http/1.1"

	if web.handler == nil {
		web.err = Errors().HandlerIsNotSet()
		return
	}
	if web.cfg == nil {
		web.err = Errors().NoConfiguration()
		return
	}
	ret = &http.Server{
		Handler:                      web.handler,
		Addr:                         web.cfg.HostPort(),
		DisableGeneralOptionsHandler: web.cfg.DisableGeneralOptionsHandler,
		ReadTimeout:                  web.cfg.ReadTimeout,
		ReadHeaderTimeout:            web.cfg.ReadHeaderTimeout,
		WriteTimeout:                 web.cfg.WriteTimeout,
		IdleTimeout:                  web.cfg.IdleTimeout,
		MaxHeaderBytes:               web.cfg.MaxHeaderBytes,
		TLSNextProto:                 make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	if web.cfg.TLSPrivateKeyPEM == "" || web.cfg.TLSPublicKeyPEM == "" {
		ret.TLSNextProto, isTLS = nil, false
		return
	}
	// Конфигурация TLS по умолчанию.
	if tlsConfig == nil {
		tlsConfig, web.err = web.net.NewTLSConfigDefault(web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM)
		if web.err != nil {
			return
		}
	}
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, http11)
	ret.TLSConfig, isTLS = tlsConfig, true

	return
}
