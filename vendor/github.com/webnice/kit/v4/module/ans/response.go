package ans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/webnice/dic"
)

// Status Ответ кодом статуса без передачи тела сообщения.
func (ans *impl) Status(wr http.ResponseWriter, status dic.IStatus) Interface {
	wr.WriteHeader(status.Code())
	return ans
}

// Ok Ответ кодом 200 "Ok".
func (ans *impl) Ok(wr http.ResponseWriter) Interface { return ans.Status(wr, dic.Status().Ok) }

// NoContent Ответ кодом 204 "No Content" без передачи тела сообщения.
func (ans *impl) NoContent(wr http.ResponseWriter) Interface {
	return ans.Status(wr, dic.Status().NoContent)
}

// Unauthorized Ответ кодом 401 "Unauthorized" без передачи тела сообщения.
func (ans *impl) Unauthorized(wr http.ResponseWriter) Interface {
	return ans.Status(wr, dic.Status().Unauthorized)
}

// Forbidden Ответ кодом 403 "Forbidden" без передачи тела сообщения.
func (ans *impl) Forbidden(wr http.ResponseWriter) Interface {
	return ans.Status(wr, dic.Status().Forbidden)
}

// NotFound Ответ кодом 404 "Not Found" без передачи тела сообщения.
func (ans *impl) NotFound(wr http.ResponseWriter) Interface {
	return ans.Status(wr, dic.Status().NotFound)
}

// BadRequest Ответ на запрос с передачей ошибки запроса и структуры описывающей найденную ошибку.
func (ans *impl) BadRequest(wr http.ResponseWriter, data RestErrorInterface) Interface {
	ans.Status(wr, dic.Status().BadRequest)
	if data == nil {
		return ans
	}
	data.Json(wr)

	return ans
}

// BadRequestBytes Ответ на запрос с передачей данных в исходном виде.
func (ans *impl) BadRequestBytes(wr http.ResponseWriter, data []byte) Interface {
	var err error

	ans.Status(wr, dic.Status().BadRequest)
	if len(data) <= 0 {
		return ans
	}
	if _, err = wr.Write(data); err != nil {
		ans.logErrorf(errResponse, err)
		return ans
	}

	return ans
}

// InternalServerError Ответ на запрос с кодом ошибки 500 и структурой описывающей ошибку.
func (ans *impl) InternalServerError(wr http.ResponseWriter, err error) Interface {
	if err != nil {
		ans.logErrorf(err.Error())
	}
	ans.
		NewRestError(dic.Status().InternalServerError, err).
		CodeSet(-1).
		Json(wr)

	return ans
}

// Header Установка заголовка передаваемых данных.
func (ans *impl) Header(wr http.ResponseWriter, header dic.IHeader, mime dic.IMime) Interface {
	return ans.HeaderString(wr, header, mime.String())
}

// HeaderString Установка заголовка передаваемых данных объектом строка.
func (ans *impl) HeaderString(wr http.ResponseWriter, header dic.IHeader, mimeString string) Interface {
	wr.Header().Set(header.String(), mimeString)
	return ans
}

// ContentType Установка типа контента передаваемых данных.
func (ans *impl) ContentType(wr http.ResponseWriter, mime dic.IMime) Interface {
	return ans.Header(wr, dic.Header().ContentType, mime)
}

// ContentTypeString Установка типа контента передаваемых данных объектом строка.
func (ans *impl) ContentTypeString(wr http.ResponseWriter, mimeString string) Interface {
	return ans.HeaderString(wr, dic.Header().ContentType, mimeString)
}

// ContentLength Установка заголовка длинны передаваемого контента.
func (ans *impl) ContentLength(wr http.ResponseWriter, contentLength uint64) Interface {
	return ans.HeaderString(wr, dic.Header().ContentLength, strconv.FormatUint(contentLength, 10))
}

// LastModified Установка заголовка с датой и временем изменения контента.
func (ans *impl) LastModified(wr http.ResponseWriter, lastModified time.Time) Interface {
	const timeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	if lastModified.IsZero() {
		return ans
	}
	return ans.HeaderString(wr, dic.Header().LastModified, lastModified.UTC().Format(timeFormat))
}

// RetryAfter Установка заголовка 'Retry-After' с числом секунд.
func (ans *impl) RetryAfter(wr http.ResponseWriter, duration time.Duration) Interface {
	return ans.HeaderString(wr, dic.Header().RetryAfter, strconv.FormatUint(uint64(duration/time.Second), 10))
}

// ResponseBytes Ответ с проверкой передачи данных.
func (ans *impl) ResponseBytes(wr http.ResponseWriter, status dic.IStatus, data []byte) Interface {
	return ans.Response(wr, status, bytes.NewBuffer(data))
}

// ResponseString Ответ с проверкой передачи данных.
func (ans *impl) ResponseString(wr http.ResponseWriter, status dic.IStatus, data string) Interface {
	return ans.Response(wr, status, bytes.NewBufferString(data))
}

// Response Ответ с проверкой передачи данных.
func (ans *impl) Response(wr http.ResponseWriter, status dic.IStatus, buf *bytes.Buffer) Interface {
	var err error

	if buf == nil {
		wr.WriteHeader(status.Code())
		return ans
	}
	if buf.Len() > 0 {
		wr.Header().Set(dic.Header().ContentLength.String(), strconv.FormatInt(int64(buf.Len()), 10))
	}
	wr.WriteHeader(status.Code())
	if buf.Len() > 0 {
		if _, err = buf.WriteTo(wr); err != nil {
			ans.logErrorf(errResponse, err)
			return ans
		}
	}

	return ans
}

// Json Ответ на запрос с сериализацией результата в JSON формат.
func (ans *impl) Json(wr http.ResponseWriter, status dic.IStatus, obj any) Interface {
	const sliceEmpty = "[]"
	var (
		err     error
		buf     *bytes.Buffer
		enc     *json.Encoder
		rvo     reflect.Value
		length  int
		isSlice bool
	)

	// Для среза получаем длину.
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		isSlice, rvo = true, reflect.ValueOf(obj)
		length = rvo.Len()
	default:
		length = 0
	}
	// Если срез пустой, отвечаем константой.
	switch {
	case isSlice && length == 0:
		buf = bytes.NewBufferString(sliceEmpty)
	default:
		buf = &bytes.Buffer{}
		enc = json.NewEncoder(buf)
		if err = enc.Encode(obj); err != nil {
			err = fmt.Errorf(errEncode, err)
			ans.InternalServerError(wr, err)
			return ans
		}
	}
	ans.ContentType(wr, dic.Mime().ApplicationJson)
	ans.Response(wr, status, buf)

	return ans
}
