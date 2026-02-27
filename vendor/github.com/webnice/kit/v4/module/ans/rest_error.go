package ans

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/webnice/dic"
)

// NewRestError Создание объекта интерфейса, стандартного REST ответа с ошибкой для кодов ошибок 4xx и 5xx.
func (ans *impl) NewRestError(status dic.IStatus, err error) (reo RestErrorInterface) {
	reo = &RestError{
		parent: ans,
		status: status,
		Error:  RestErrorBody{Code: -1},
	}
	if err != nil {
		reo.MessageSet(err.Error())
	}

	return
}

// CodeSet Уникальный код ошибки из справочника ошибок.
func (reo *RestError) CodeSet(code int) RestErrorInterface { reo.Error.Code = code; return reo }

// CodeGet Уникальный код ошибки из справочника ошибок.
func (reo *RestError) CodeGet() int { return reo.Error.Code }

// MessageSet Техническое описание ошибки в локализации сервера.
func (reo *RestError) MessageSet(msg string) RestErrorInterface { reo.Error.Message = msg; return reo }

// MessageGet Техническое описание ошибки в локализации сервера.
func (reo *RestError) MessageGet() string { return reo.Error.Message }

// I18nKeySet Ключ локализации ошибки.
func (reo *RestError) I18nKeySet(key string) RestErrorInterface { reo.Error.I18nKey = key; return reo }

// I18nKeyGet Ключ локализации ошибки.
func (reo *RestError) I18nKeyGet() string { return reo.Error.I18nKey }

// Field Массив объектов с описанием имён полей и ошибок в них.
func (reo *RestError) Field() []RestErrorField { return reo.Error.Errors }

// Add Добавление в массив объектов описания поля и ошибки в нём.
func (reo *RestError) Add(field string, value string, msg string) RestErrorInterface {
	reo.Error.Errors = append(reo.Error.Errors, RestErrorField{
		Field:      field,
		FieldValue: value,
		Message:    msg,
	})

	return reo
}

// AddWithKey Добавление в массив объектов описания поля и ошибки в нём.
func (reo *RestError) AddWithKey(field string, value string, msg string, key string) RestErrorInterface {
	reo.Error.Errors = append(reo.Error.Errors, RestErrorField{
		Field:      field,
		FieldValue: value,
		Message:    msg,
		I18nKey:    key,
	})

	return reo
}

// JsonBytes Сериализация в JSON.
func (reo *RestError) JsonBytes() (ret *bytes.Buffer, err error) {
	var enc *json.Encoder

	ret = &bytes.Buffer{}
	enc = json.NewEncoder(ret)
	if err = enc.Encode(reo); err != nil {
		ret, err = nil, fmt.Errorf(errEncode, err)
		return
	}

	return
}

// Json Сериализация в JSON и отправка данных в интерфейс http.ResponseWriter.
func (reo *RestError) Json(wr http.ResponseWriter) {
	var (
		err error
		buf *bytes.Buffer
	)

	if buf, err = reo.JsonBytes(); err != nil {
		reo.parent.InternalServerError(wr, err)
		return
	}
	reo.parent.ContentType(wr, dic.Mime().ApplicationJson)
	reo.parent.Response(wr, reo.status, buf)

	return
}

// AsError Сводит объект в ошибку и возвращает интерфейс error.
func (reo *RestError) AsError() (err error) {
	const (
		errTemplate  = "ошибки значений в полях структуры данных"
		errCodeError = "ошибка с кодом %d"
	)

	if reo.Error.Message != "" {
		err = errors.New(reo.Error.Message)
		return
	}
	if len(reo.Error.Errors) > 0 {
		err = errors.New(errTemplate)
		return
	}
	if reo.Error.Code != 0 {
		err = fmt.Errorf(errCodeError, reo.Error.Code)
		return
	}

	return
}
