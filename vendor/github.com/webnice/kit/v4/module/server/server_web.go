package server

import (
	"sync"

	kitModuleUuid "github.com/webnice/kit/v4/module/uuid"
	kitTypes "github.com/webnice/kit/v4/types"
	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// Создание объекта и возвращение интерфейса InterfaceWeb.
func newWeb(parent *impl) *implWeb {
	var iweb = &implWeb{
		parent: parent,
		lib:    newWebLib(parent),
		error:  newWebError(parent),
		server: &implWebServer{
			Protect: new(sync.RWMutex),
			Control: make(map[string]*kitTypesServer.WebServerControl),
		},
	}

	return iweb
}

// Ссылка на менеджер логирования.
func (iweb *implWeb) log() kitTypes.Logger { return iweb.parent.log() }

// Отправка сообщения в лог уровня "info", если включён режим отладки.
func (iweb *implWeb) info(pattern string, args ...any) {
	if iweb.parent.debug {
		iweb.log().Infof(pattern, args...)
	}
}

// Add Добавление конфигурации веб сервера. Функцию можно вызывать многократно, для добавления
// нескольких веб серверов. Серверы не должны пересекаться по занимаемому IP и порту или другим монопольно
// выделяемым ресурсам.
// Возвращается объект добавленного веб сервера.
func (iweb *implWeb) Add(cfg *kitTypesServer.WebConfiguration) (ret *kitTypesServer.Server) {
	ret = &kitTypesServer.Server{
		T:   kitTypesServer.TWeb,
		Web: cfg,
	}
	// Создание идентификатора сервера.
	if ret.Web.Server.ID == "" {
		ret.Web.Server.ID = kitModuleUuid.Get().V4().String()
	}
	// Добавление конфигурации веб сервера.
	iweb.parent.server = append(iweb.parent.server, ret)

	return
}

// Reg Регистрация ресурсов сервера.
func (iweb *implWeb) Reg(res kitTypesServer.WebResource) { iweb.res = append(iweb.res, res) }

// Lib Библиотека ВЕБ функций.
func (iweb *implWeb) Lib() InterfaceWebLib { return iweb.lib }

// Error Интерфейс ошибок ВЕБ сервера.
func (iweb *implWeb) Error() InterfaceWebError { return iweb.error }
