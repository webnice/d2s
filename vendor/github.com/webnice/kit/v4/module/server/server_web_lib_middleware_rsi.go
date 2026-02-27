package server

import (
	"context"
	"net"
	"net/http"

	kitModuleRqVar "github.com/webnice/kit/v4/module/rqvar"
	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// RequestShadowInfoHandler Загрузка объекта информации о запросе и всех данных переданных в запросе в контекст.
func (iwl *implWebLib) RequestShadowInfoHandler() (ret func(http.Handler) http.Handler) {
	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			rq = rq.WithContext(context.WithValue(rq.Context(), cKeyRequestShadowInfo, &contextWrapper{
				Name:  cNameRequestShadowInfo,
				Value: iwl.RequestShadowInfoExtraction(rq),
			}))
			next.ServeHTTP(wr, rq)
		})
		return
	}

	return
}

// RequestShadowInfoExtraction Получение информации о запросе и всех данных переданных в запросе.
func (iwl *implWebLib) RequestShadowInfoExtraction(rq *http.Request) (ret *kitTypesServer.RequestShadowInfo) {
	var (
		err error
		ipo net.IP
		vsc *kitTypesServer.WebServerControl
	)

	ret = new(kitTypesServer.RequestShadowInfo)
	if vsc, err = iwl.parent.Web().Lib().Middleware().WebServerControlGetFromContext(rq); err == nil {
		ret.Server = vsc.Configuration
	}
	if err = kitModuleRqVar.Get().Load(rq, ret); err != nil {
		ret = nil
		return
	}
	if ipo = iwl.parent.Web().Lib().Middleware().IpGetFromContext(rq); len(ipo) > 0 {
		ret.IP = ipo
	}

	return
}

// RequestShadowInfoGetFromContext Извлечение объекта информации о запросе и всех данных переданных в запросе
// из контекста HTTP запроса.
// Возвращается nil, если объект информации не был найден в контексте HTTP запроса.
func (iwl *implWebLib) RequestShadowInfoGetFromContext(rq *http.Request) (ret *kitTypesServer.RequestShadowInfo) {
	var (
		ok      bool
		wrapper *contextWrapper
	)

	if wrapper = iwl.getContextWrapper(rq, cKeyRequestShadowInfo, cNameRequestShadowInfo); wrapper == nil {
		return
	}
	if ret, ok = wrapper.Value.(*kitTypesServer.RequestShadowInfo); !ok {
		return
	}

	return
}
