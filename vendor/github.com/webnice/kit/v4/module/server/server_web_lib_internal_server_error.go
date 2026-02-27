package server

import (
	"bytes"
	"net/http"

	"github.com/webnice/dic"
	"github.com/webnice/kit/v4/module/ans"
)

// InternalServerErrorGet Получение функции для обработки внутренней ошибки ВЕБ сервера.
func (iwl *implWebLib) InternalServerErrorGet() (ret http.HandlerFunc) {
	if ret = iwl.funcs.fnInternalServerError; ret == nil {
		ret = iwl.defaultInternalServerError
	}

	return
}

// InternalServerErrorSet Установка пользовательской функции для обработки внутренней ошибки ВЕБ сервера.
func (iwl *implWebLib) InternalServerErrorSet(fn http.HandlerFunc) InterfaceHandlerFunc {
	if fn != nil {
		iwl.funcs.fnInternalServerError = fn
	}

	return iwl
}

func (iwl *implWebLib) defaultInternalServerError(wr http.ResponseWriter, _ *http.Request) {
	var (
		ero error
		rsp ans.Interface
		buf *bytes.Buffer
	)

	rsp = ans.New(iwl.parent.logger)
	rsp.ContentType(wr, dic.Mime().TextPlain)
	buf = bytes.NewBuffer(dic.Status().InternalServerError.Bytes())
	if ero = iwl.parent.serverWeb.error.InternalServerError(nil); ero != nil {
		buf = bytes.NewBufferString(ero.Error())
		iwl.parent.log().Error(ero.Error())
	}
	rsp.ResponseBytes(wr, dic.Status().InternalServerError, buf.Bytes())
}
