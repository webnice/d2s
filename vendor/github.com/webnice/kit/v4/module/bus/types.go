package bus

/*

Шина данных.

*/

import (
	"context"
	"reflect"

	kitModulePdw "github.com/webnice/kit/v4/module/pdw"
	kitTypes "github.com/webnice/kit/v4/types"
)

const (
	defaultDatabusBufferLength = 1000000 // Размер буферизированного канала шины данных по умолчанию.
	defaultWorkerCount         = 1000    // Количество параллельно работающих обработчиков, обслуживающих канал данных.
)

// Объект сущности пакета.
type impl struct {
	debug         bool               // Флаг отладки.
	workerCount   int64              // Счётчик выполняющихся обработчиков шины данных.
	workerContext context.Context    // Контекст остановки обработчиков шины данных.
	workerCancel  context.CancelFunc // Функция прерывания работы обработчиков шины данных.
	workerOnDone  chan struct{}      // Канал подтверждения остановки обработчика шины данных.
	databus       *databus           // Внутренние объекты на которых реализуется шина данных.
	essence       Essence            // Объект интерфейса Essence.
}

// Объект сути сущности, интерфейс Essence.
type gist struct {
	parent *impl // Адрес объекта основной сущности, интерфейс Interface.
}

// Внутренние объекты на которых реализуется шина данных.
type databus struct {
	Wrappers    kitModulePdw.Interface // Бассейн объектов для обёртки передаваемых потребителю данных.
	Subscribers *subscribers           // Подписчики.
	Bus         chan kitModulePdw.Data // Канал обмена данными.
}

// Описание подписчика
type subscriber struct {
	Name  string                // Подписчик: пакет, объект.
	Item  kitTypes.Databuser    // Интерфейс подписчика.
	Types []*subscriberTypeInfo // Типы объектов, ожидаемые потребителем.
}

// Описание типа данных, которыми ожидают подписчики.
type subscriberTypeInfo struct {
	Original         any          // Оригинальный объект.
	OriginalType     reflect.Type // Тип данных в том виде как был передан в регистрацию.
	TypeName         string       // Название типа.
	BaseType         reflect.Type // Базовый тип после раскрытия всех адресов.
	IndirectionCount int          // Количество раскрытий адресов до получения базового типа.
}

// Структура с результатами обработки данных.
type workerSafeCallResponse struct {
	Err  error   // Ошибка вызова потребителя.
	Resp []any   // Результат обработки данных потребителем.
	Errs []error // Ошибки возвращённые потребителем.
}
