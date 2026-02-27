package cfg

import (
	"time"

	kitTypes "github.com/webnice/kit/v4/types"
)

// Recorder Интерфейс логирования.
type Recorder kitTypes.Logger

// Структура объекта сущности логирования, интерфейс Logger.
type recorder struct {
	parent *impl // Адрес объекта основной сущности, интерфейс Interface.
}

// Структура объекта переноса данных.
type record struct {
	recorder         *recorder           // Родительский объект.
	timestamp        time.Time           // Время записи в UTC.
	traceInfo        *kitTypes.TraceInfo // Данные о вызывающей функции и стеке.
	stackBack        int                 // Дополнительный откат записей стека.
	stackBackCorrect int                 // Коррекция откат записей стека для отдельно взятого вызова логирования.
	keys             map[string]any      // Ключи логирования.
	fatality         *bool               // Флаг фатальности записи лога.
}
