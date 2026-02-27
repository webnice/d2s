package server

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/webnice/dic"
)

// IpHandler Загрузка IP адреса клиента в контекст.
func (iwl *implWebLib) IpHandler() (ret func(http.Handler) http.Handler) {
	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			rq = rq.WithContext(context.WithValue(rq.Context(), cKeyIP, &contextWrapper{
				Name:  cNameIP,
				Value: iwl.IpExtraction(rq),
			}))
			next.ServeHTTP(wr, rq)
		})
		return
	}

	return
}

// IpExtraction Получение IP адреса клиента запроса.
func (iwl *implWebLib) IpExtraction(rq *http.Request) (ret net.IP) {
	const headerXClientForwardedFor = "X-Client-Forwarded-For"
	var (
		err error
		tmp string
		i   int
	)

	// Приоритетный заголовок
	if tmp = rq.Header.Get(headerXClientForwardedFor); tmp != "" {
		ret = net.ParseIP(tmp)
		return
	}
	if tmp = strings.TrimSpace(rq.Header.Get(dic.Header().XRealIP.String())); tmp != "" {
		ret = net.ParseIP(tmp)
		return
	}
	tmp = rq.Header.Get(dic.Header().XForwardedFor.String())
	if i = strings.IndexByte(tmp, ','); i >= 0 {
		tmp = tmp[:i]
	}
	if tmp = strings.TrimSpace(tmp); tmp != "" {
		ret = net.ParseIP(tmp)
		return
	}
	if tmp = strings.TrimSpace(rq.RemoteAddr); tmp == "" {
		return
	}
	if tmp, _, err = net.SplitHostPort(tmp); err != nil {
		return
	}
	ret = net.ParseIP(tmp)

	return
}

// IpGetFromContext Извлечение объекта IP адреса клиента из контекста HTTP запроса.
// Возвращается nil - когда http.Handler IP адреса не был подключен в "промежуточный слой".
func (iwl *implWebLib) IpGetFromContext(rq *http.Request) (ret net.IP) {
	var (
		ok      bool
		wrapper *contextWrapper
	)

	if wrapper = iwl.getContextWrapper(rq, cKeyIP, cNameIP); wrapper == nil {
		return
	}
	if ret, ok = wrapper.Value.(net.IP); !ok {
		return
	}

	return
}
