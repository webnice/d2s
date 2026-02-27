package ans

import (
	"bytes"
	"net/http"
)

// RestErrorInterface Интерфейс стандартного REST ответа ошибки с кодами 4xx и 5xx.
type RestErrorInterface interface {
	// CodeSet Уникальный код ошибки из справочника ошибок.
	CodeSet(code int) RestErrorInterface

	// CodeGet Уникальный код ошибки из справочника ошибок.
	CodeGet() int

	// MessageSet Назначение значения полю message.
	MessageSet(msg string) RestErrorInterface

	// MessageGet Возвращение значения поля message.
	MessageGet() string

	// I18nKeySet Ключ локализации ошибки.
	I18nKeySet(key string) RestErrorInterface

	// I18nKeyGet Ключ локализации ошибки.
	I18nKeyGet() string

	// Field Массив объектов с описанием имён полей и ошибок в них.
	Field() []RestErrorField

	// Add Добавление в массив объектов описания поля и ошибки в нём.
	Add(field string, value string, msg string) RestErrorInterface

	// AddWithKey Добавление в массив объектов описания поля и ошибки в нём.
	AddWithKey(field string, value string, msg string, key string) RestErrorInterface

	// JsonBytes Сериализация в JSON.
	JsonBytes() (ret *bytes.Buffer, err error)

	// Json Сериализация в JSON и отправка данных в интерфейс http.ResponseWriter.
	Json(wr http.ResponseWriter)

	// AsError Сводит объект в ошибку и возвращает интерфейс error.
	AsError() (err error)
}
