package server

import (
	"sync"

	kitTypesServer "github.com/webnice/kit/v4/types/server"

	"github.com/go-chi/chi/v5"
)

// InterfaceWeb Интерфейс веб сервера.
type InterfaceWeb interface {
	// Add Добавление конфигурации веб сервера. Функцию можно вызывать многократно, для добавления
	// нескольких веб серверов. Серверы не должны пересекаться по занимаемому IP и порту или другим монопольно
	// выделяемым ресурсам.
	// Возвращается объект добавленного веб сервера.
	Add(cfg *kitTypesServer.WebConfiguration) (ret *kitTypesServer.Server)

	// Reg Регистрация ресурсов сервера.
	Reg(res kitTypesServer.WebResource)

	// Lib Библиотека ВЕБ функций.
	Lib() InterfaceWebLib

	// Error Интерфейс ошибок ВЕБ сервера.
	Error() InterfaceWebError

	// IsStarted Функция вернёт булево значение "истина", если ВЕБ сервер уже запущен.
	IsStarted(serverID string) (ret bool)

	// Start Запуск ВЕБ сервера с указанным идентификатором.
	Start(serverID string) (err error)

	// Stop Остановка ВЕБ сервера с указанным идентификатором.
	Stop(serverID string) (err error)
}

// Объект сущности, реализующий интерфейс InterfaceWeb.
type implWeb struct {
	parent *impl                        // Адрес объекта родительской сущности (parent), интерфейс Interface.
	lib    *implWebLib                  // Объект сущности интерфейса InterfaceWebLib.
	error  *implWebError                // Объект сущности интерфейса InterfaceWebError.
	server *implWebServer               // Объекты запущенных ВЕБ серверов.
	res    []kitTypesServer.WebResource // Зарегистрированные ресурсы ВЕБ сервера.
	router *chi.Mux                     // Роутер веб сервера.
}

// Объекты запущенных ВЕБ серверов.
type implWebServer struct {
	Protect *sync.RWMutex
	Control map[string]*kitTypesServer.WebServerControl
}
