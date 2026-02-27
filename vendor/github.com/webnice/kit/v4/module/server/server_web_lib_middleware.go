package server

import "net/http"

const (
	// Название контекста для объекта контроля ВЕБ сервера.
	cNameWebServerControl = "middleware-web-server-control"

	// Название контекста для объекта ip адреса клиента HTTP запроса.
	cNameIP = "middleware-client-request-ip"

	// Название контекста для объекта получения полной информации о запросе и всех данных переданных в запросе.
	cNameRequestShadowInfo = "middleware-request-shadow-info"

	// Название контекста для объекта простой web авторизации.
	cNameBasicAuth = "middleware-request-basic-auth"
)

var (
	// Ключ контекста для объекта контроля ВЕБ сервера.
	cKeyWebServerControl = &contextKey{cNameWebServerControl}

	// Ключ контекста для объекта ip адреса клиента HTTP запроса.
	cKeyIP = &contextKey{cNameIP}

	// Ключ контекста для объекта получения полной информации о запросе и всех данных переданных в запросе.
	cKeyRequestShadowInfo = &contextKey{cNameRequestShadowInfo}

	// Ключ контекста для объекта простой web авторизации.
	cKeyBasicAuth = &contextKey{cNameBasicAuth}
)

// Общая функция извлечения из контекста обёртки над данными с проверкой ключа и названия.
func (iwl *implWebLib) getContextWrapper(rq *http.Request, key any, name string) (ret *contextWrapper) {
	var (
		ok    bool
		value any
	)

	if value = rq.Context().Value(key); value == nil {
		return
	}
	if ret, ok = value.(*contextWrapper); !ok {
		return
	}
	if ret == nil || ret.Value == nil || ret.Name != name {
		ret = nil
		return
	}

	return
}
