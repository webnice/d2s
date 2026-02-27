package server

import (
	"context"
	"fmt"
	"net/http"

	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// WebServerControlHandler Обработчик установки в контекст запросов ВЕБ сервера, объекта контроля за ВЕБ сервером.
func (iwl *implWebLib) WebServerControlHandler(ctl *kitTypesServer.WebServerControl) (
	ret func(http.Handler) http.Handler,
) {
	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			rq = rq.WithContext(
				context.WithValue(rq.Context(), cKeyWebServerControl, &contextWrapper{
					Name:  cNameWebServerControl,
					Value: ctl,
				}))
			next.ServeHTTP(wr, rq)
		})
		return
	}

	return
}

// WebServerControlGetFromContext Извлечение объекта контроля за ВЕБ сервером из контекста HTTP запроса.
func (iwl *implWebLib) WebServerControlGetFromContext(rq *http.Request) (
	ret *kitTypesServer.WebServerControl,
	err error,
) {
	const (
		tplNotFound   = "данные по ключу %q в контексте не найдены"
		tplWrongValue = "по ключу %q получена обёртка контекста содержащая не верные данные"
	)
	var (
		ok      bool
		wrapper *contextWrapper
	)

	if wrapper = iwl.getContextWrapper(rq, cKeyWebServerControl, cNameWebServerControl); wrapper == nil {
		err = fmt.Errorf(tplNotFound, cNameWebServerControl)
		return
	}
	if ret, ok = wrapper.Value.(*kitTypesServer.WebServerControl); !ok {
		err = fmt.Errorf(tplWrongValue, cNameWebServerControl)
		return
	}

	return
}
