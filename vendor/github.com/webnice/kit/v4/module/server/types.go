/*

Сервер реализует обработку доступа к ресурсам с использованием: ВЕБ, REST, GRPC, TCP, UDP.

*/

package server

import (
	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitTypes "github.com/webnice/kit/v4/types"
	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

const moduleName = "module/server"

// Interface Интерфейс пакета.
type Interface interface {
	// Gist Интерфейс служебных методов.
	Gist() Essence

	// Web Интерфейс веб сервера.
	Web() InterfaceWeb

	// Grpc Интерфейс GRPC сервера.
	//Grpc() kitTypesServer.IServer[kitTypesServer.Grpc]

	// Tcp Интерфейс TCP/IP сервера.
	//Tcp() kitTypesServer.IServer[kitTypesServer.Tcp]

	// Udp Интерфейс UDP сервера.
	//Udp() kitTypesServer.IServer[kitTypesServer.Udp]

	// Start Запуск всех зарегистрированных серверов.
	// Если не зарегистрирован ни один сервер, функция возвращает ошибку.
	Start() (err error)

	// Stop Остановка всех зарегистрированных серверов.
	// Если не зарегистрирован ни один сервер, функция возвращает ошибку.
	Stop() (err error)

	// Errors Справочник ошибок.
	Errors() *Error
}

// Essence Служебный публичный интерфейс.
type Essence interface {
	// Debug Присвоение нового значения режима отладки.
	Debug(debug bool) Essence
}

// Объект сущности пакета.
type impl struct {
	debug     bool                     // Флаг режима отладки.
	logger    kitTypes.Logger          // Интерфейс менеджера логирования.
	gist      *gist                    // Объект служебного интерфейса Essence.
	server    []*kitTypesServer.Server // Добавленные конфигурации серверов.
	serverWeb *implWeb                 // Объект служебного интерфейса InterfaceWeb.
	answer    kitModuleAns.Interface   // Интерфейс работы с http запросами и ответами.
}

// Объект сути сущности, интерфейс Essence.
type gist struct {
	parent *impl // Адрес объекта родительской сущности (parent), интерфейс Interface.
}
