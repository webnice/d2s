package web

import (
	"crypto/tls"
	"net"
	"reflect"
	"slices"
	"strings"
	"time"
)

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных как ":http".
//
//goland:noinspection DuplicatedCode
func parseAddress(addr string) (ret *Configuration, err error) {
	const bColon = ":"
	var (
		sp []string
		n  int
	)

	addr = strings.TrimSpace(addr)
	ret, sp = new(Configuration), make([]string, 2)
	if sp[0], sp[1], err = net.SplitHostPort(addr); err != nil {
		ret.Host, err = addr, nil
		return
	}
	if n, err = net.LookupPort(netTcp, strings.Join(sp[1:], bColon)); err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint16(n)
	switch sp, err = net.LookupHost(sp[0]); err {
	case nil:
		ret.Host = sp[0]
	default:
		err = nil
	}

	return
}

// http2RequiredCipherSuites Шифры, обязательные для работы протокола HTTP/2.
var http2RequiredCipherSuites = []uint16{
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
}

// ensureHTTP2CipherSuites Проверка списка шифров TLS на наличие хотя бы одного шифра, обязательного для
// протокола HTTP/2. Если ни одного обязательного шифра в списке нет, они добавляются.
func ensureHTTP2CipherSuites(cipherSuites []uint16) (ret []uint16) {
	ret = cipherSuites
	if len(ret) == 0 {
		// Пустой список шифров означает использование списка по умолчанию, он совместим с HTTP/2.
		return
	}
	for _, cipher := range http2RequiredCipherSuites {
		if slices.Contains(cipherSuites, cipher) {
			return
		}
	}
	ret = append(ret, http2RequiredCipherSuites...)

	return
}

// isTLSListener Проверка того, что слушатель соединений уже обёрнут в TLS режим.
// Тип слушателя TLS из пакета "crypto/tls" не экспортируется, поэтому проверка выполняется по пакету типа.
func isTLSListener(ltn net.Listener) (ok bool) {
	const pkgTls = "crypto/tls"
	var tp = reflect.TypeOf(ltn)

	if tp == nil {
		return
	}
	for tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	ok = tp.PkgPath() == pkgTls

	return
}

// Заполнение конфигурации значениями по умолчанию.
func defaultConfiguration(cfg *Configuration) {
	const defaultShutdownTimeout = time.Duration(30) * time.Second

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = defaultShutdownTimeout
	}
}
