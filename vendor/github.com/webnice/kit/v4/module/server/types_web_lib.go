package server

import (
	"net"
	"net/http"

	kitTypesServer "github.com/webnice/kit/v4/types/server"
)

// InterfaceWebLib Интерфейс библиотеки функций ВЕБ сервера.
type InterfaceWebLib interface {
	// БИБЛИОТЕКИ

	// Handler Интерфейс библиотеки ВЕБ функций Handler и HandlerFunc.
	Handler() InterfaceHandlerFunc

	// Middleware Интерфейс библиотеки ВЕБ функций "промежуточного слоя" для роутинга на базе библиотеки chi.
	Middleware() InterfaceMiddleware
}

// InterfaceHandlerFunc Интерфейс библиотеки ВЕБ функций Handler и HandlerFunc.
type InterfaceHandlerFunc interface {
	// ФУНКЦИИ

	// InternalServerErrorGet Получение функции для обработки внутренней ошибки ВЕБ сервера.
	InternalServerErrorGet() (ret http.HandlerFunc)

	// InternalServerErrorSet Установка пользовательской функции для обработки внутренней ошибки ВЕБ сервера.
	InternalServerErrorSet(fn http.HandlerFunc) InterfaceHandlerFunc
}

// InterfaceMiddleware Интерфейс библиотеки ВЕБ функций "промежуточного слоя".
type InterfaceMiddleware interface {
	// ПРОМЕЖУТОЧНЫЙ СЛОЙ

	// IpHandler Загрузка IP адреса клиента в контекст.
	IpHandler() (ret func(http.Handler) http.Handler)

	// IpGetFromContext Извлечение объекта IP адреса клиента из контекста HTTP запроса.
	// Возвращается nil - когда http.Handler IP адреса не был подключен в "промежуточный слой".
	IpGetFromContext(rq *http.Request) net.IP

	// LogHandler Запись в журнал запросов к ВЕБ серверу.
	LogHandler() (ret func(http.Handler) http.Handler)

	// RecoverHandler Обработчик восстановления после паники в ВЕБ сервере.
	RecoverHandler() (ret func(http.Handler) http.Handler)

	// WebServerControlHandler Обработчик установки в контекст запросов ВЕБ сервера, объекта контроля за ВЕБ сервером.
	WebServerControlHandler(ctl *kitTypesServer.WebServerControl) (ret func(http.Handler) http.Handler)

	// WebServerControlGetFromContext Извлечение объекта контроля за ВЕБ сервером из контекста HTTP запроса.
	WebServerControlGetFromContext(rq *http.Request) (ret *kitTypesServer.WebServerControl, err error)

	// RequestShadowInfoHandler Загрузка объекта информации о запросе и всех данных переданных в запросе в контекст.
	RequestShadowInfoHandler() (ret func(http.Handler) http.Handler)

	// RequestShadowInfoGetFromContext Извлечение объекта информации о запросе и всех данных переданных в запросе
	// из контекста HTTP запроса.
	// Возвращается nil, если объект информации не был найден в контексте HTTP запроса.
	RequestShadowInfoGetFromContext(rq *http.Request) (ret *kitTypesServer.RequestShadowInfo)

	// BasicAuthHandler Простая web авторизация запросов с использованием логина и пароля.
	BasicAuthHandler(cfg kitTypesServer.BasicAuthConfiguration) (ret func(http.Handler) http.Handler)

	// BasicAuthGetFromContext Извлечение имени пользователя из контекста проверки авторизации через простую авторизацию.
	BasicAuthGetFromContext(rq *http.Request) (ret string)
}

// Объект сущности, реализующий интерфейс InterfaceWebLib.
type implWebLib struct {
	parent *impl      // Адрес объекта родительской сущности (parent), интерфейс Interface.
	funcs  *fnHandler // Функции HandlerFunc.
}

// Функции Handler и HandlerFunc.
type fnHandler struct {
	fnInternalServerError http.HandlerFunc // Функция обработки ошибки InternalServerError.
}

// Обёртка для формирования ключа данных контекста ВЕБ запросов.
type contextKey struct {
	name string
}

// Обёртка над данными размещаемыми в контексте ВЕБ запросов.
type contextWrapper struct {
	Name  string
	Value any
}
