package server

import (
	"runtime"
	"strings"

	kitTypesServer "github.com/webnice/kit/v4/types/server"
	"github.com/webnice/web/v3"

	"github.com/go-chi/chi/v5"
)

// IsStarted Функция вернёт булево значение "истина", если ВЕБ сервер уже запущен.
func (iweb *implWeb) IsStarted(serverID string) (ret bool) {
	var (
		ok  bool
		ctl *kitTypesServer.WebServerControl
	)

	// Блокировка конкурентного доступа.
	iweb.server.Protect.RLock()
	defer iweb.server.Protect.RUnlock()
	serverID = strings.ToUpper(serverID)
	if ctl, ok = iweb.server.Control[serverID]; !ok {
		return
	}
	ret = ctl != nil

	return
}

// Start Запуск ВЕБ сервера с указанным идентификатором.
func (iweb *implWeb) Start(serverID string) (err error) {
	const patternBeg, patternEnd = "Запуск ВЕБ сервера %q.", "Запуск ВЕБ сервера %q, выполнен."
	var (
		serverKey string
		ok        bool
		n         int
		router    *chi.Mux
	)

	// Блокировка конкурентного доступа.
	iweb.server.Protect.Lock()
	defer iweb.server.Protect.Unlock()
	// Проверка наличия сервера в списке.
	serverKey = strings.ToUpper(serverID)
	if _, ok = iweb.server.Control[serverKey]; ok {
		err = iweb.parent.Errors().ServerByIdAlreadyStarted.Bind(serverKey)
		return
	}
	iweb.server.Control[serverKey] = new(kitTypesServer.WebServerControl)
	// Поиск конфигурации сервера по ID.
	for n = range iweb.parent.server {
		if strings.EqualFold(iweb.parent.server[n].Web.Server.ID, serverKey) {
			iweb.server.Control[serverKey].Configuration = iweb.parent.server[n]
		}
	}
	if iweb.server.Control[serverKey].Configuration == nil {
		err = iweb.parent.Errors().ServerByIdNotFound.Bind(serverKey)
		return
	}
	iweb.info(patternBeg, serverID)
	iweb.server.Control[serverKey].Server = web.
		New()
	// Обработчик добавления в контекст объекта ВЕБ сервера, конфигурация и интерфейс контроля.
	router = chi.NewRouter()
	router.Use(iweb.Lib().Middleware().WebServerControlHandler(iweb.server.Control[serverKey]))
	router.Mount("/", iweb.router)
	iweb.server.Control[serverKey].Server.Handler(router)
	// Запуск ВЕБ сервера.
	if err = iweb.server.Control[serverKey].
		Server.
		ListenAndServeWithConfig(&iweb.server.Control[serverKey].Configuration.Web.Server).
		Error(); err != nil {
		return
	}
	iweb.info(patternEnd, iweb.server.Control[serverKey].Server.ID())

	return
}

// Stop Остановка ВЕБ сервера с указанным идентификатором.
func (iweb *implWeb) Stop(serverID string) (err error) {
	const (
		patternBeg, patternWit = "Остановка ВЕБ сервера %q.", "Ожидание остановки ВЕБ сервера %q."
		patternEnd             = "Остановка ВЕБ сервера %q, выполнена."
	)
	var (
		serverKey string
		ok        bool
		si        web.Interface
	)

	runtime.Gosched()
	iweb.info(patternBeg, serverID)
	// Блокировка конкурентного доступа.
	iweb.server.Protect.Lock()
	defer iweb.server.Protect.Unlock()
	serverKey = strings.ToUpper(serverID)
	// Проверка наличия сервера в списке.
	if _, ok = iweb.server.Control[serverKey]; !ok {
		err = iweb.parent.Errors().ServerByIdNotStarted.Bind(serverKey)
		return
	}
	// Остановка сервера, ожидание завершения сервера.
	iweb.info(patternWit, serverID)
	// DEBUG
	//println("Stop.начало.")
	// DEBUG
	si = iweb.server.Control[serverKey].
		Server.
		Stop()
	// DEBUG
	//println("Stop.окончание")
	// DEBUG
	runtime.Gosched()
	err = si.Wait().Error()
	// Удаление сервера из списка запущенных.
	delete(iweb.server.Control, serverKey)
	iweb.info(patternEnd, serverID)

	return
}
