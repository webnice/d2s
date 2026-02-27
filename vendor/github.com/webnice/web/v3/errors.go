package web

import wnet "github.com/webnice/net"

// Все ошибки определены как константы.
const (
	cHandlerIsNotSet = "Не установлен обработчик запросов ВЕБ сервера."
)

// Константы указываются в объектах в качестве фиксированного адреса на протяжении всего времени работы приложения.
// Ошибка с ошибкой могут сравниваться по содержимому, по адресу и т.д.
var (
	errSingleton       = &Error{}
	errHandlerIsNotSet = err(cHandlerIsNotSet)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }

// ОШИБКИ.

// AlreadyRunning Сервер уже запущен.
func (e *Error) AlreadyRunning() error { return wnet.Errors().AlreadyRunning() }

// NoConfiguration Конфигурация сервера отсутствует либо равна nil.
func (e *Error) NoConfiguration() error { return wnet.Errors().NoConfiguration() }

// ListenSystemdPID Переменная окружения LISTEN_PID пустая, либо содержит не верное значение.
func (e *Error) ListenSystemdPID() error { return wnet.Errors().ListenSystemdPID() }

// ListenSystemdFDS Переменная окружения LISTEN_FDS пустая, либо содержит не верное значение.
func (e *Error) ListenSystemdFDS() error { return wnet.Errors().ListenSystemdFDS() }

// ListenSystemdNotFound Получение сокета systemd по имени, имя не найдено.
func (e *Error) ListenSystemdNotFound() error { return wnet.Errors().ListenSystemdNotFound() }

// ListenSystemdQuantityNotMatch Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES.
func (e *Error) ListenSystemdQuantityNotMatch() error {
	return wnet.Errors().ListenSystemdQuantityNotMatch()
}

// TLSIsNil Конфигурация TLS сервера пустая.
func (e *Error) TLSIsNil() error { return wnet.Errors().TLSIsNil() }

// ServerHandlerIsNotSet Не установлен обработчик основной функции TCP сервера.
func (e *Error) ServerHandlerIsNotSet() error { return wnet.Errors().ServerHandlerIsNotSet() }

// HandlerIsNotSet Не установлен обработчик запросов ВЕБ сервера.
func (e *Error) HandlerIsNotSet() error { return &errHandlerIsNotSet }
