package ans

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/webnice/dic"
	kitModuleRqVar "github.com/webnice/kit/v4/module/rqvar"

	"github.com/go-chi/chi/v5"
)

// RqVar Интерфейс модуля github.com/webnice/kit/module/rqvar. Поиск и загрузка данных из:
// + контекста HTTP запроса;
// + заголовков HTTP запроса;
// + "печенек" HTTP запроса;
// + параметров HTTP запроса, а именно в URL Request;
// + путей URN роутинга HTTP запроса;
// + вызов функции структуры загружаемых данных;
func (ans *impl) RqVar() kitModuleRqVar.Interface { return kitModuleRqVar.Get() }

// RqIds Загрузка числовых идентификаторов указанных в path-param под ключём {key} и перечисленных через запятую.
// В случае возникновения ошибки формируется и отправляется HTTP ответ с кодом 400, содержащий возникшую ошибку.
func (ans *impl) RqIds(wr http.ResponseWriter, rq *http.Request, key string) (ret []uint64, err error) {
	const sep = ","
	var (
		srcIDs string
		sIDs   []string
		n      int
		v      uint64
		tmp    string
	)

	srcIDs = chi.URLParam(rq, key)
	sIDs = strings.Split(srcIDs, sep)
	ret = make([]uint64, 0, len(sIDs))
	for n = range sIDs {
		if tmp = strings.TrimSpace(sIDs[n]); tmp == "" {
			continue
		}
		if v, err = strconv.ParseUint(tmp, 10, 64); err != nil {
			err = fmt.Errorf("конвертация значения %q в число прервана ошибкой: %s", tmp, err)
			ans.NewRestError(dic.Status().BadRequest, err).CodeSet(-1).
				Add(key, tmp, err.Error()).Json(wr)
			return
		}
		ret = append(ret, v)
	}

	return
}

// RqId Загрузка значения указанного в path-param под ключём {key} и конвертация его в число.
// В случае возникновения ошибки формируется и отправляется HTTP ответ с кодом 400, содержащий возникшую ошибку.
func (ans *impl) RqId(wr http.ResponseWriter, rq *http.Request, key string) (ret uint64, err error) {
	var (
		src string
	)

	if src = strings.TrimSpace(chi.URLParam(rq, key)); src == "" {
		return
	}
	if ret, err = strconv.ParseUint(src, 10, 64); err != nil {
		err = fmt.Errorf("конвертация значения %q в число прервана ошибкой: %s", src, err)
		ans.NewRestError(dic.Status().BadRequest, err).CodeSet(-1).
			Add(key, src, err.Error()).Json(wr)
		return
	}

	return
}

// RqLoad Выполнение загрузки и декодирования данных, в случае возникновения ошибки формируется и отправляется HTTP
// ответ с кодом 400, содержащий возникшую ошибку.
// Данные загружаются из тела запроса в переменную variable с использованием декодирования данных, выбор кодера
// осуществляется на основе заголовка запроса Content-Type, поддерживаются два метода сериализации данных: JSON, XML.
// Ответ с ошибкой сериализуется тем же самым методом сериализации данных, что был в запросе.
func (ans *impl) RqLoad(wr http.ResponseWriter, rq *http.Request, variable any) (err error) {
	if err = ans.RqData(rq, variable); err != nil {
		ans.NewRestError(dic.Status().BadRequest, err).Json(wr)
		return
	}

	return
}

// RqLoadVerify Выполнение загрузки и декодирования данных, в случае возникновения ошибки формируется и
// отправляется HTTP ответ с кодом 400, содержащий возникшую ошибку.
// Загруженные данные проверяются библиотекой github.com/go-playground/validator.
// Данные загружаются из тела запроса в переменную variable с использованием декодирования данных, выбор кодера
// осуществляется на основе заголовка запроса Content-Type, поддерживаются два метода сериализации данных: JSON, XML.
// Ответ с ошибкой сериализуется тем же самым методом сериализации данных, что был в запросе.
func (ans *impl) RqLoadVerify(wr http.ResponseWriter, rq *http.Request, variable any) (err error) {
	var reo RestErrorInterface

	if err = ans.RqLoad(wr, rq, variable); err != nil {
		return
	}
	if reo = ans.Verify(variable); reo == nil {
		return
	}
	err = reo.AsError()
	reo.Json(wr)

	return
}
