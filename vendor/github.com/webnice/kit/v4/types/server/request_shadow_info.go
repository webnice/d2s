package server

import (
	"net"
	"net/http"
	"net/url"
)

// RequestShadowInfo Информация извлекаемая из запроса.
type RequestShadowInfo struct {
	Server                     *Server
	ProjectID                  uint64  `header:"X-Project-Id"`                                                                // Уникальный идентификатор проекта.
	Authorization              string  `header:"Authorization"   cookie:"Access-Token"   urn-param:"accessToken,accesstoken"` // Текущий сессионный токен доступа.
	UserAgent                  string  `header:"User-Agent"`                                                                  // Заголовок, значение которого используется для описания устройства пользователя.
	AcceptLanguage             string  `header:"Accept-Language"`                                                             // Предпочитаемые языки локализации.
	AcceptEncoding             string  `header:"Accept-Encoding"`                                                             // Перечень поддерживаемых способов кодирования для ответа на запрос.
	DeviceID                   string  `header:"X-Device-Id"     cookie:"AuthDeviceID"`                                       // Содержит идентификатор ранее сохранённого устройства пользователя.
	DeviceName                 string  `header:"X-Device-Name"`                                                               // Название устройства или название приложения, для отображения в списке запомненных устройств.
	IP                         net.IP  `header:"X-Real-Ip,X-Client-Forwarded-For"`                                            // IP адрес клиента.
	Scheme                     string  `header:"X-Client-Forwarded-Scheme"`                                                   // Протокол по которому подключается клиент, возможные значения: http, https.
	Domain                     string  `header:"X-Client-Forwarded-Domain,Host"          call-func:"RequestExtractionDomain"` // Доменное имя сервера к которому пришел запрос.
	Origin                     string  `header:"Origin"`                                                                      // Источник, схема, хост, порт запроса.
	ContentType                string  `header:"Content-Type"`                                                                // Формат передаваемых данных.
	Singleton                  bool    `                                                 urn-param:"singleton"`               // Режим одиночки. Если указано true, при успешной аутентификации все другие действующие сессии пользователя будут удалены.
	NoCookie                   bool    `                                                 urn-param:"noCookie"`                // Не использовать Cookie.
	LocationLatitude           float64 `header:"X-Location-Latitude"`                                                         // Координаты широты местонахождения устройства.
	LocationLongitude          float64 `header:"X-Location-Longitude"`                                                        // Координаты долготы местонахождения устройства.
	LocationAltitude           float64 `header:"X-Location-Altitude"`                                                         // Высота местонахождения устройства над уровнем моря в метрах.
	LocationSpeed              float64 `header:"X-Location-Speed"`                                                            // Скорость передвижения устройства в метрах в секунду.
	LocationAzimuth            float64 `header:"X-Location-Azimuth"`                                                          // Азимут (направление движения) в градусах относительно истинного севера.
	LocationAccuracyHorizontal float64 `header:"X-Location-Accuracy-Horizontal"`                                              // Радиус неопределённости местоположения по горизонтали в метрах.
	LocationAccuracyVertical   float64 `header:"X-Location-Accuracy-Vertical"`                                                // Радиус неопределённости местоположения по вертикали в метрах.
	LocationAccuracySpeed      float64 `header:"X-Location-Accuracy-Speed"`                                                   // Точность значения скорости, измеряемое в метрах в секунду.
	LocationAccuracyAzimuth    float64 `header:"X-Location-Accuracy-Azimuth"`                                                 // Точность измерения направления движения устройства в градусах.
	LocationTimestamp          uint64  `header:"X-Location-Timestamp"`                                                        // Дата и время полученное со спутников в момент определения координат.
}

// Поиск и извлечение имени домена из конфигурации сервера, если конфигурация передана.
func (rsi *RequestShadowInfo) extractionDomainFromServer(s *Server) (ret string) {
	var (
		err error
		uri *url.URL
	)

	if s == nil {
		return
	}
	if s.T != TWeb {
		return
	}
	if s.Web == nil {
		return
	}
	if ret = s.Web.Server.Address; ret != "" {
		uri, err = url.Parse(ret)
		if ret = ""; err == nil {
			ret = uri.Host
			return
		}
	}
	if len(s.Web.Server.Domain) > 0 {
		if ret = s.Web.Server.Domain[0]; ret != "" {
			return
		}
	}
	ret = s.Web.Server.Host

	return
}

func (rsi *RequestShadowInfo) RequestExtractionDomain(rq *http.Request) (ret string) {
	// Если домен уже был найден другим способом - выход.
	if ret = rsi.Domain; ret != "" {
		return
	}
	if ret = rq.URL.Hostname(); ret != "" {
		return
	}
	if ret = rsi.extractionDomainFromServer(rsi.Server); ret != "" {
		return
	}

	return
}
