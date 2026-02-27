package server

import (
	"fmt"
	"net/http"

	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// RecoverHandler Восстановление после паники в ВЕБ сервере.
func (iwl *implWebLib) RecoverHandler() (ret func(http.Handler) http.Handler) {
	const tpl = "Восстановление после паники:\n%v\n%s."

	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			defer func() {
				var e any
				if e = recover(); e != nil {
					// Установка значения ошибки.
					_ = iwl.parent.serverWeb.error.InternalServerError(fmt.Errorf(tpl, e, kitModuleTrace.StackShort()))
					// Вызов функции-обработчика InternalServerError().
					iwl.InternalServerErrorGet().ServeHTTP(wr, rq)
				}
			}()
			next.ServeHTTP(wr, rq)
		})
		return
	}

	return
}
