package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/webnice/dic"
	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// BasicAuthHandler Простая web авторизация запросов с использованием логина и пароля.
func (iwl *implWebLib) BasicAuthHandler(cfg kitTypesServer.BasicAuthConfiguration) (
	ret func(http.Handler) http.Handler,
) {
	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			var isNext bool
			if rq, isNext = iwl.basicAuth(cfg, wr, rq); !isNext {
				return
			}
			next.ServeHTTP(wr, rq)
		})
		return
	}

	return
}

// Выполнение проверки простой web авторизации.
func (iwl *implWebLib) basicAuth(cfg kitTypesServer.BasicAuthConfiguration, wr http.ResponseWriter, rq *http.Request) (
	ret *http.Request,
	isNext bool,
) {
	const (
		charset  = "UTF-8"
		tplRealm = "Basic realm=%q, charset=%q"
		defaultQ = "Пожалуйста, авторизуйтесь."
		errCall  = "выполнение функции проверки простой авторизации прервано ошибкой: %s"
	)
	var (
		err      error
		rsp      kitTypesServer.BasicAuthResponse
		username string
		password string
		realmQry string
		ok       bool
		hKey     dic.IHeader
	)

	// Получение данных запроса.
	ret = rq
	username, password, ok = rq.BasicAuth()
	// Запрос авторизации.
	if !ok {
		if realmQry = defaultQ; cfg.Request != "" {
			realmQry = cfg.Request
		}
		realmQry = fmt.Sprintf(tplRealm, realmQry, charset)
		iwl.parent.ans().HeaderString(wr, dic.Header().WwwAuthenticate, realmQry)
		if cfg.Header != nil {
			for hKey = range cfg.Header {
				iwl.parent.ans().HeaderString(wr, hKey, cfg.Header[hKey])
			}
		}
		switch cfg.Body {
		case nil:
			iwl.parent.ans().ResponseBytes(wr, dic.Status().Unauthorized, dic.Status().Unauthorized.Bytes())
		default:
			iwl.parent.ans().ResponseBytes(wr, dic.Status().Unauthorized, cfg.Body.Bytes())
		}
		return
	}
	// Выполнение проверки полученных данных авторизации.
	if err = iwl.safeCall(func() { rsp = cfg.AuthFunc(username, password) }); err != nil {
		err = fmt.Errorf(errCall, err)
		iwl.parent.ans().InternalServerError(wr, err)
		return
	}
	// Предоставлены не верные данные.
	if !rsp.IsCorrect {
		if rsp.Header != nil {
			for hKey = range rsp.Header {
				iwl.parent.ans().HeaderString(wr, hKey, rsp.Header[hKey])
			}
		}
		switch rsp.Body {
		case nil:
			switch rsp.Code {
			case nil:
				iwl.parent.ans().ResponseBytes(wr, dic.Status().Unauthorized, dic.Status().Unauthorized.Bytes())
			default:
				iwl.parent.ans().ResponseBytes(wr, rsp.Code, rsp.Code.Bytes())
			}
		default:
			iwl.parent.ans().ResponseBytes(wr, rsp.Code, rsp.Body.Bytes())
		}
		return
	}
	// Предоставлены верные данные, сохранение данных в контекст.
	ret = rq.WithContext(context.WithValue(rq.Context(), cKeyBasicAuth, &contextWrapper{
		Name:  cNameBasicAuth,
		Value: username,
	}))
	isNext = true

	return
}

// BasicAuthGetFromContext Извлечение имени пользователя из контекста проверки авторизации через простую авторизацию.
func (iwl *implWebLib) BasicAuthGetFromContext(rq *http.Request) (ret string) {
	var (
		ok      bool
		wrapper *contextWrapper
	)

	if wrapper = iwl.getContextWrapper(rq, cKeyBasicAuth, cNameBasicAuth); wrapper == nil {
		return
	}
	if ret, ok = wrapper.Value.(string); !ok {
		return
	}

	return
}
