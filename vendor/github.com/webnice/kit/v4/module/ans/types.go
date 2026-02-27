package ans

import (
	"github.com/webnice/dic"
	kitTypes "github.com/webnice/kit/v4/types"
)

const (
	errEncode   = "сериализация в JSON прервана ошибкой: %s"
	errResponse = "передача HTTP ответа прервана ошибкой: %s"
)

// Essence Служебный публичный интерфейс.
type Essence interface {
	// Debug Присвоение нового значения режима отладки.
	Debug(debug bool) Essence
}

// Объект сущности пакета.
type impl struct {
	debug   bool            // Флаг режима отладки.
	essence Essence         // Объект интерфейса Essence.
	logger  kitTypes.Logger // Интерфейс менеджера логирования.
}

// Объект сути сущности, интерфейс Essence.
type gist struct {
	parent *impl // Адрес объекта родительской сущности (parent), интерфейс Interface.
}

// RestError Структура стандартного REST ответа с ошибкой с кодами 4xx и 5xx.
type RestError struct {
	parent *impl
	status dic.IStatus
	Error  RestErrorBody `json:"error" xml:"error"`
}

// RestErrorBody Общее описание ошибки.
type RestErrorBody struct {
	Code    int              `json:"code"    xml:"code"`              // Уникальный код ошибки из справочника ошибок.
	Message string           `json:"message" xml:"message"`           // Техническое описание ошибки в локализации сервера.
	I18nKey string           `json:"i18nKey" xml:"i18nKey,omitempty"` // Ключ локализации ошибки.
	Errors  []RestErrorField `json:"errors"  xml:"errors"`            // Массив объектов с описанием имён полей и ошибок в них.
}

// RestErrorField Описание ошибок в полях структуры данных.
type RestErrorField struct {
	Field      string `json:"field"      xml:"field"`             // Название поля.
	FieldValue string `json:"fieldValue" xml:"fieldValue"`        // Полученное значение поля field.
	Message    string `json:"message"    xml:"message"`           // Техническое описание ошибки.
	I18nKey    string `json:"i18nKey"    xml:"i18NKey,omitempty"` // Ключ локализации ошибки.
}
