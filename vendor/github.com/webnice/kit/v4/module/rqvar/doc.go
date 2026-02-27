package rqvar

/*

Модуль загрузки данных в структуру из http.Request и параметров и путей URN роутинга.
Какие данные откуда брать описывается в тегах структуры данных.

Для поиска, извлечения и загрузки данных, доступны следующие теги:
  - "context"    - В стандартном контексте HTTP запроса;
  - "header"     - В заголовках HTTP запроса;
  - "cookie"     - В "печеньках" HTTP запроса;
  - "urn-param"  - В параметрах HTTP запроса, а именно в URL Request;
  - "path-param" - В пути URN роутинга, пример: /api/v1.0/path/{ArticleID}/main.
                   Где 'ArticleID' - это параметр пути, при использовании роутера github.com/go-chi/chi;
  - "call-func"  - Через вызов функции структуры с указанным именем, тип функции описан в "rqvar.RqFunc";

Все теги могут использоваться одновременно, в каждом из тегов могут быть перечислены большее одного ключа поиска данных.
Ключи не могут быть пустыми, равняться значению '-' и должны быть перечислены через запятую.
Когда указано множество тегов и множество ключей, библиотека ищет данные последовательно, в порядке перечисления ключей
и до появления первых не пустых данных, которые записываются в поле объекта структуры.
Приоритет последовательности поиска данных следующий:
  1: context
  2: header
  3: cookie
  4: urn-param
  5: path-param
  6: call-func

Функция структуры получения данных, помимо извлечения или добычи данных, может их обрабатывать, а так же может собой
заменить все другие теги, так как код функции может быть любой сложности. Это сделано для универсальности и гибкости.

Пример структуры:
type ExampleVariableStruct struct {
	Authorization              string  `header:"Authorization" cookie:"Access-Token" urn-param:"accessToken,accesstoken"`
	UserAgent                  string  `header:"User-Agent"`
	AcceptLanguage             string  `header:"Accept-Language"`
	DeviceID                   string  `header:"X-Device-Id"      cookie:"AuthDeviceID"`
	IP                         net.IP  `header:"X-Client-Forwarded-For,X-Real-IP,X-Forwarded-For"`
    LocationLatitude           float64 `call-func:"LoadLocationLatitude"`
    ArticleID                  uint64  `path-param:"ArticleID"`
    SessionID                  string  `context:"AuthorizationSessionID"`
	Err                        error   `header:"-"   cookie:"-"   urn-param:"-"   path-param:"-"   call-func:"-"`
}

// LoadLocationLatitude Функция получения данных для поля структуры LocationLatitude.
func (evs *ExampleVariableStruct) LoadLocationLatitude(rq *http.Request) string {...}


Типы данных объекта структуры, в которые библиотека может присваивать найденные значения:
  - Все базовые типы данных языка go;
  - net.IP;


Если объект структуры содержит поля с вложенными структурами, тогда все вложенные структуры игнорируются.

*/
