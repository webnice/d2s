package server

import (
	"fmt"

	kitModuleAns "github.com/webnice/kit/v4/module/ans"
	kitTypes "github.com/webnice/kit/v4/types"
	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New(logger kitTypes.Logger) Interface {
	var sri = &impl{
		logger: logger,
		answer: kitModuleAns.New(logger),
	}
	sri.gist = newEssence(sri)
	sri.serverWeb = newWeb(sri)

	return sri
}

// Ссылка на менеджер логирования.
func (sri *impl) log() kitTypes.Logger { return sri.logger }

// Интерфейс библиотеки функций работы с HTTP запросами и ответами.
func (sri *impl) ans() kitModuleAns.Interface { return sri.answer }

// Errors Справочник ошибок.
func (sri *impl) Errors() *Error { return Errors() }

// Gist Интерфейс служебных методов.
func (sri *impl) Gist() Essence { return sri.gist }

// Web Интерфейс веб сервера.
func (sri *impl) Web() InterfaceWeb { return sri.serverWeb }

// Grpc Интерфейс GRPC сервера.
//func (sri *impl) Grpc() InterfaceGrpc { return sri.serverGrpc }

// Tcp Интерфейс TCP сервера.
//func (sri *impl) Tcp() InterfaceTcp { return sri.serverTcp }

// Udp Интерфейс UDP сервера.
//func (sri *impl) Udp() InterfaceUdp { return sri.serverUdp }

// Start Запуск всех зарегистрированных серверов.
func (sri *impl) Start() (err error) {
	var (
		serverAdded    uint64
		serverLaunched uint64
		n              int
	)

	serverAdded = uint64(len(sri.server))
	// Подготовка зарегистрированных ресурсов веб сервера.
	for n = range sri.server {
		switch sri.server[n].T {
		case kitTypesServer.TWeb:
			err = sri.serverWeb.Prepare(sri.server[n].Web)
		}
		if err != nil {
			err = fmt.Errorf("подготовка ресурсов сервера прервана ошибкой: %s", err)
			return
		}
	}
	// Подготовка зарегистрированных ресурсов GRPC сервера(ов).
	//if err = sri.serverGrpc.
	//	Prepare(); err != nil {
	//	err = fmt.Errorf("подготовка ресурсов сервера прервана ошибкой: %s", err)
	//	return
	//}
	// Подготовка зарегистрированных ресурсов TCP сервера(ов).
	//if err = sri.serverTcp.
	//	Prepare(); err != nil {
	//	err = fmt.Errorf("подготовка ресурсов сервера прервана ошибкой: %s", err)
	//	return
	//}
	// Подготовка зарегистрированных ресурсов UDP сервера(ов).
	//if err = sri.serverUdp.
	//	Prepare(); err != nil {
	//	err = fmt.Errorf("подготовка ресурсов сервера прервана ошибкой: %s", err)
	//	return
	//}
	// Запуск всех зарегистрированных серверов.
	for n = range sri.server {
		switch sri.server[n].T {
		case kitTypesServer.TWeb:
			err = sri.serverWeb.
				Start(sri.server[n].Web.Server.ID)
		}
		switch err {
		case nil:
			serverLaunched++
		default:
			sri.log().
				Criticalf("запуск %q сервера прерван ошибкой: %s", sri.server[n].T.String(), err)
			err = nil
		}
	}
	// Количество запущенных серверов не равно количеству добавленных серверов.
	if serverAdded != serverLaunched {
		err = sri.Errors().ServersAddedNotEqualLaunched.Bind(serverAdded, serverLaunched)
		return
	}

	return
}

// Stop Остановка всех зарегистрированных серверов.
func (sri *impl) Stop() (err error) {
	var n int

	// Остановка всех зарегистрированных серверов.
	for n = range sri.server {
		switch sri.server[n].T {
		case kitTypesServer.TWeb:
			err = sri.serverWeb.
				Stop(sri.server[n].Web.Server.ID)
		}
		if err != nil {
			sri.log().
				Warningf("остановка %q сервера завершилась ошибкой: %s", sri.server[n].T.String(), err)
			err = nil
		}
	}

	return
}
