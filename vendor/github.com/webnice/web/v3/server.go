package web

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
)

// maxBodyBytesReader Обёртка тела запроса, ограничивающая его размер и отвечающая ошибкой
// 413 Request Entity Too Large при превышении лимита.
type maxBodyBytesReader struct {
	w http.ResponseWriter
	r io.ReadCloser
}

// Read Чтение тела запроса.
func (mbr *maxBodyBytesReader) Read(p []byte) (n int, err error) {
	if n, err = mbr.r.Read(p); err != nil {
		var errMax *http.MaxBytesError
		if //goland:noinspection GoFixErrorsAsType
		errors.As(err, &errMax) {
			// Лимит размера тела запроса превышен. Ответ 413 возможен если обработчик ещё не начал ответ.
			http.Error(mbr.w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		}
	}

	return
}

// Close Закрытие тела запроса.
func (mbr *maxBodyBytesReader) Close() error { return mbr.r.Close() }

// maxBodyBytesHandler Ограничение максимального размера тела запроса в n байт с ответом ошибкой
// 413 Request Entity Too Large при превышении лимита.
func maxBodyBytesHandler(h http.Handler, n int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var r2 = new(http.Request)

		*r2 = *r
		r2.Body = &maxBodyBytesReader{w: w, r: http.MaxBytesReader(w, r.Body, n)}
		h.ServeHTTP(w, r2)
	})
}

// Применение настроек из конфигурации к стандартному веб серверу.
func (web *impl) makeServer(tlsConfig *tls.Config) (ret *http.Server, isTLS bool) {
	const (
		http11 = "http/1.1"
		proto2 = "h2"
	)
	var handler http.Handler

	if web.handler == nil {
		web.err = Errors().HandlerIsNotSet()
		return
	}
	if web.cfg == nil {
		web.err = Errors().NoConfiguration()
		return
	}
	// Ограничение максимального размера тела запроса.
	handler = web.handler
	if n := int64(web.cfg.MaxBodyBytes); n > 0 {
		handler = maxBodyBytesHandler(handler, n)
	}
	ret = &http.Server{
		Handler:                      handler,
		Addr:                         web.cfg.HostPort(),
		DisableGeneralOptionsHandler: web.cfg.DisableGeneralOptionsHandler,
		ReadTimeout:                  web.cfg.ReadTimeout,
		ReadHeaderTimeout:            web.cfg.ReadHeaderTimeout,
		WriteTimeout:                 web.cfg.WriteTimeout,
		IdleTimeout:                  web.cfg.IdleTimeout,
		MaxHeaderBytes:               web.cfg.MaxHeaderBytes,
	}
	// Настройка протоколов HTTP. HTTP/1 всегда включён.
	ret.Protocols = new(http.Protocols)
	ret.Protocols.SetHTTP1(true)
	if !web.cfg.HTTP2Disable {
		ret.Protocols.SetHTTP2(true)
	}
	if web.cfg.HTTP2UnencryptedEnable {
		ret.Protocols.SetUnencryptedHTTP2(true)
	}
	// Настройка протокола HTTP/2.
	if ret.Protocols.HTTP2() || ret.Protocols.UnencryptedHTTP2() {
		ret.HTTP2 = &http.HTTP2Config{
			MaxConcurrentStreams: web.cfg.HTTP2MaxConcurrentStreams,
		}
	}
	if web.cfg.TLSPrivateKeyPEM == "" || web.cfg.TLSPublicKeyPEM == "" {
		isTLS = false
		return
	}
	// Конфигурация TLS по умолчанию.
	if tlsConfig == nil {
		tlsConfig, web.err = web.net.NewTLSConfigDefault(web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM)
		if web.err != nil {
			return
		}
	}
	// Порядок протоколов определяет предпочтение сервера при согласовании ALPN.
	if ret.Protocols.HTTP2() {
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, proto2)
		// Сервер с включённым протоколом HTTP/2 не запускается, если в списке шифров отсутствует хотя бы
		// один шифр, обязательный для HTTP/2.
		tlsConfig.CipherSuites = ensureHTTP2CipherSuites(tlsConfig.CipherSuites)
	}
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, http11)
	ret.TLSConfig, isTLS = tlsConfig, true

	return
}
