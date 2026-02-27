package ans

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/webnice/dic"
)

// RqBytes Загрузка тела HTTP запроса в виде среза байт и возвращение объекта *bytes.Buffer.
func (ans *impl) RqBytes(rq *http.Request) (ret *bytes.Buffer, err error) {
	const (
		errHttpRequest = "в качестве адреса *http.Request был передан nil"
		errIoCopy      = "чтение тела HTTP запроса прервано ошибкой: %s"
		errLength      = "размер тела HTTP запроса (байт): %d, загружено (байт): %d"
	)
	var copied int64

	if rq == nil {
		err = errors.New(errHttpRequest)
		return
	}
	defer func() {
		if e := rq.Body.Close(); e != nil {
			ans.logWarningf("Закрытие тела HTTP запроса, после чтения данных, прервано ошибкой: %s.", e)
			return
		}
	}()
	// Получение тела запроса.
	ret = &bytes.Buffer{}
	if copied, err = io.Copy(ret, rq.Body); err != nil {
		err = fmt.Errorf(errIoCopy, err)
		return
	}
	if rq.ContentLength > 0 {
		if copied != rq.ContentLength {
			err = fmt.Errorf(errLength, rq.ContentLength, copied)
			return
		}
	}

	return
}

// RqData Выполнение загрузки данных из тела запроса в переменную variable с использованием декодирования
// данных, выбор кодера осуществляется на основе заголовка запроса Content-Type.
// Поддерживаются два метода сериализации данных: JSON, XML.
func (ans *impl) RqData(rq *http.Request, variable any) (err error) {
	const (
		errEmpty              = "тело запроса пустое"
		errUnknownContentType = "заголовок Content-Type не передан, либо неизвестен тип контента: %q"
		errDecoderError       = "ошибка декодирования: %s"
	)
	var (
		buf *bytes.Buffer
		cti dic.IMime
	)

	// Загрузка тела HTTP запроса в виде среза байт.
	if buf, err = ans.RqBytes(rq); err != nil {
		return
	}
	if buf.Len() < 2 {
		err = errors.New(errEmpty)
		return
	}
	// Выбор типа кодирования на основе Content-Type заголовка запроса.
	cti = dic.ParseMime(rq.Header.Get(dic.Header().ContentType.String()))
	switch {
	case dic.Mime().ApplicationJson.IsEqual(cti):
		err = json.NewDecoder(buf).Decode(variable)
	case dic.Mime().TextXml.IsEqual(cti), dic.Mime().ApplicationXml.IsEqual(cti):
		err = xml.NewDecoder(buf).Decode(variable)
	default:
		err = fmt.Errorf(errUnknownContentType, rq.Header.Get(dic.Header().ContentType.String()))
		return
	}
	// Ошибка декодирования.
	if err != nil {
		err = fmt.Errorf(errDecoderError, err)
	}

	return
}
