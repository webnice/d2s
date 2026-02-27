package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/webnice/dic"
	kitModuleDye "github.com/webnice/kit/v4/module/dye"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitModuleWrWrap "github.com/webnice/kit/v4/module/wrwrap"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Record Структура собираемых данных для логирования.
type Record struct {
	Address  net.IP        // IP адрес клиента.
	Code     int           // Код ответа.
	Method   dic.IMethod   // Метод запроса.
	Size     uint64        // Размер ответа в байтах.
	Path     string        // Запрашиваемый путь.
	LeadTime time.Duration // Время выполнения запроса.
}

// LogHandler Запись в журнал запросов к ВЕБ серверу.
func (iwl *implWebLib) LogHandler() (ret func(http.Handler) http.Handler) {
	ret = func(next http.Handler) (fn http.Handler) {
		fn = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			var (
				beginTime time.Time                 // Дата и время начала обработки запроса.
				mwrw      kitModuleWrWrap.Interface // Обёртка над http.ResponseWriter.
			)

			beginTime = time.Now()
			mwrw = kitModuleWrWrap.New(wr, iwl.parent.debug)
			defer iwl.logHandler(beginTime, mwrw, rq)
			next.ServeHTTP(mwrw, rq)
		})
		return
	}

	return
}

// Формирование строки лога и печать в лог.
func (iwl *implWebLib) logHandler(beginTime time.Time, wrw kitModuleWrWrap.Interface, rq *http.Request) {
	var (
		logRecord    Record
		logLevel     kitModuleLogLevel.Level
		log          kitTypes.Logger
		codeColorBeg string
		codeColorEnd string
	)

	// Режим "продакшн" - пишем в журнал только коды ошибок больше 499.
	if !iwl.parent.debug && wrw.StatusCode() <= 499 {
		return
	}
	// Вычисление времени выполнения запроса.
	logRecord.LeadTime = time.Since(beginTime)
	// Загрузка IP адреса клиента HTTP запроса.
	logRecord.Address = iwl.Middleware().IpGetFromContext(rq)
	// Сбор данных запроса.
	logRecord.Method = dic.ParseMethod(rq.Method)
	logRecord.Path = rq.URL.Path
	logRecord.Code = wrw.StatusCode()
	logRecord.Size = wrw.Len()
	// Если используется http.ServeContent() или http.ServeFile().
	if logRecord.Size == 0 {
		logRecord.Size, _ = strconv.ParseUint(wrw.Header().Get(dic.Header().ContentLength.String()), 10, 64)
	}
	// Определение цвета записи в журнале в зависимости от кода HTTP ответа.
	codeColorEnd = kitModuleDye.New().Reset().Done().String()
	switch {
	case logRecord.Code >= 100 && logRecord.Code < 200:
		codeColorBeg = kitModuleDye.New().Underline().Bright().Yellow().Done().String()
	case logRecord.Code >= 200 && logRecord.Code < 300:
		codeColorBeg = kitModuleDye.New().Underline().Bright().White().Done().String()
	case logRecord.Code >= 300 && logRecord.Code < 400:
		codeColorBeg = kitModuleDye.New().Underline().Bright().Cyan().Done().String()
	case logRecord.Code >= 400 && logRecord.Code < 500:
		codeColorBeg = kitModuleDye.New().Underline().Bright().Magenta().Done().String()
	case logRecord.Code >= 500 && logRecord.Code <= 599:
		codeColorBeg = kitModuleDye.New().Underline().Bright().Red().Done().String()
	default:
		codeColorBeg = kitModuleDye.New().Underline().Normal().Done().String()
	}
	// Проецирование уровня логирования в зависимости от кода HTTP ответа.
	switch {
	case logRecord.Code >= 500:
		logLevel = kitModuleLogLevel.Error
	case logRecord.Code >= 400:
		logLevel = kitModuleLogLevel.Warning
	default:
		logLevel = kitModuleLogLevel.Info
	}
	// Режим "продакшн" - игнорируются все запросы со штатным HTTP кодом ответа.
	if !iwl.parent.debug && logLevel.Int() >= 200 && logLevel.Int() < 300 {
		return
	}
	// Запись в журнал с разбивкой по полям, для возможности работы с GrayLog.
	log = iwl.parent.log().
		Key(
			kitTypes.LoggerKey{"Address": logRecord.Address.String()},
			kitTypes.LoggerKey{"Code": logRecord.Code},
			kitTypes.LoggerKey{"Method": logRecord.Method},
			kitTypes.LoggerKey{"Size": logRecord.Size},
			kitTypes.LoggerKey{"Path": logRecord.Path},
			kitTypes.LoggerKey{"LeadTime": logRecord.LeadTime.String()},
			kitTypes.LoggerKey{"LeadTimeMillisecond": logRecord.LeadTime.Nanoseconds() / int64(time.Millisecond)},
		)
	if wrw.Header().Get(dic.Header().Location.String()) != "" {
		log = log.Key(
			kitTypes.LoggerKey{"Location": wrw.Header().Get(dic.Header().Location.String())},
		)
	}
	log.MessageWithLevel(
		logLevel,
		"%s%3d%s %s %s",
		codeColorBeg, logRecord.Code, codeColorEnd,
		logRecord.Method,
		logRecord.Path,
	)
}
