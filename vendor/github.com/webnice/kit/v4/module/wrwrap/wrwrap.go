package wrwrap

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
// responseWriter - Объект интерфейса http.ResponseWriter.
// isDebug        - Флаг режима отладки.
func New(responseWriter http.ResponseWriter, isDebug bool) Interface {
	var wrw = &impl{
		debug:   isDebug,
		essence: responseWriter,
	}

	return wrw
}

// StatusCode Статуса HTTP ответа.
func (wrw *impl) StatusCode() int {
	if wrw.statusCode == nil {
		return http.StatusOK
	}
	return *wrw.statusCode
}

// Len Длинна HTTP ответа в байтах.
func (wrw *impl) Len() uint64 { return wrw.responseLength }

// Tee Получение интерфейса io.Writer, в который передаётся копия всех данных, передаваемых в функцию Write().
func (wrw *impl) Tee(tee io.Writer) Interface { wrw.tee = tee; return wrw }

// Essence Возвращает оригинальный объект http.ResponseWriter.
func (wrw *impl) Essence() http.ResponseWriter { return wrw.essence }

// Header Возвращает адрес объекта заголовков, которые будет отправлены клиенту.
func (wrw *impl) Header() http.Header { return wrw.essence.Header() }

// Write Передаёт данные в соединение с клиентом.
func (wrw *impl) Write(bytes []byte) (n int, err error) {
	wrw.WriteHeader(http.StatusOK)
	if n, err = wrw.essence.Write(bytes); err != nil {
		return
	}
	if wrw.tee != nil {
		_, err = wrw.tee.Write(bytes[:n])
	}
	wrw.responseLength += uint64(n)

	return
}

// WriteHeader Передаёт заголовок HTTP ответа.
func (wrw *impl) WriteHeader(statusCode int) {
	if wrw.statusCode != nil {
		return
	}
	wrw.statusCode = new(int)
	*wrw.statusCode = statusCode
	wrw.essence.WriteHeader(statusCode)
}

// Flush Реализация интерфейса http.Flusher.
func (wrw *impl) Flush() {
	if wfi, ok := wrw.essence.(http.Flusher); ok {
		wfi.Flush()
	}
}

// Hijack Реализация интерфейса http.Hijacker.
func (wrw *impl) Hijack() (con net.Conn, buf *bufio.ReadWriter, err error) {
	var (
		ok  bool
		hji http.Hijacker
	)

	if hji, ok = wrw.essence.(http.Hijacker); !ok {
		err = http.ErrHijacked
		return
	}
	con, buf, err = hji.Hijack()

	return nil, nil, nil
}

// Push Реализация интерфейса http.Pusher.
func (wrw *impl) Push(target string, opts *http.PushOptions) (err error) {
	var (
		ok bool
		hp http.Pusher
	)

	if hp, ok = wrw.essence.(http.Pusher); !ok {
		err = http.ErrNotSupported
		return
	}
	err = hp.Push(target, opts)

	return
}

// WriteString Реализация интерфейса io.StringWriter.
func (wrw *impl) WriteString(s string) (n int, err error) {
	var (
		ok bool
		sw io.StringWriter
	)

	wrw.WriteHeader(http.StatusOK)
	if sw, ok = wrw.essence.(io.StringWriter); !ok {
		return wrw.Write([]byte(s))
	}
	if n, err = sw.WriteString(s); err != nil {
		return
	}
	if wrw.tee != nil {
		_, err = wrw.tee.Write([]byte(s[:n]))
	}
	wrw.responseLength += uint64(n)

	return
}
