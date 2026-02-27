package wrwrap

/*

Обёртка над http.ResponseWriter для получения статистики о данных HTTP ответа для формирования
записи для журнала логирования запросов.

*/

import (
	"io"
	"net/http"
)

// Interface Интерфейс пакета.
type Interface interface {
	http.ResponseWriter
	http.Hijacker
	http.Pusher
	http.Flusher
	io.StringWriter

	// StatusCode Статуса HTTP ответа.
	StatusCode() int

	// Len Длинна HTTP ответа в байтах.
	Len() uint64

	// Tee Получение интерфейса io.Writer, в который передаётся копия всех данных, передаваемых в функцию Write().
	Tee(io.Writer) Interface

	// Essence Возвращает оригинальный объект http.ResponseWriter.
	Essence() http.ResponseWriter
}

// Объект сущности пакета.
type impl struct {
	debug          bool                // Флаг режима отладки.
	essence        http.ResponseWriter // Оригинальный объект интерфейса http.ResponseWriter.
	statusCode     *int                // Последний код статуса HTTP ответа.
	responseLength uint64              // Количество переданных клиенту байт.
	tee            io.Writer           // Объект записи копии данных.
}
