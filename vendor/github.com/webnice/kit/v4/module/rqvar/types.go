package rqvar

import "net/http"

const (
	tagContext  = "context"    // Поиск данных в стандартном контексте HTTP запроса.
	tagHeader   = "header"     // Поиск данных в заголовках.
	tagCookie   = "cookie"     // Поиск данных в печеньках.
	tagParam    = "urn-param"  // Поиск данных в параметрах запроса.
	tagUrnParam = "path-param" // Поиск данных в пути URN роутинга.
	tagRqFunc   = "call-func"  // Поиск значений через вызов функции объекта структуры.
)

const (
	keySkip        = "-" // Значения ключа означающее пропуск поиска данных.
	separatorComma = "," // Разделитель ключей.
)

// Библиотека размещается в памяти как пакет-одиночка.
var singleton *impl

var _ = Get()

// Interface Интерфейс пакета.
type Interface interface {
	// Load Загрузка данных из запроса в объект структуры.
	Load(request *http.Request, variable any) (err error)
}

// Объект сущности пакета.
type impl struct {
}

// RqFunc Описание типа функции, доступной для вызова при поиске данных.
type RqFunc func(r *http.Request) string
