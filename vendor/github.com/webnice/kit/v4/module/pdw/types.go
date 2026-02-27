package pdw

/*

Бассейн обёрток передачи данных.

*/

import (
	"context"
	"reflect"
	"sync"
)

// Interface Интерфейс пакета.
type Interface interface {
	// Debug Присвоение нового значения режима отладки.
	Debug(debug bool) Interface

	// Get Получение объекта из бассейна.
	Get() Data

	// Put Возвращение объекта в бассейн.
	Put(wdo Data)

	// Statistic Статистика работы бассейна.
	// Статистика ведётся только если бассейн создан с флагом отладки New(isDebug=true).
	// Если бассейн создан без флага отладки, статистика вернёт nil.
	Statistic() *Statistic
}

// Data Интерфейс обёртки данных.
type Data interface {
	// DataPut Загрузка в объект данных, флага и контекста.
	DataPut(data any, isSync bool, ctx context.Context)

	// DataGet Возвращение оборачиваемых данных.
	DataGet() any

	// Type Тип обёрнутых данных.
	Type() reflect.Type

	// IsSync Флаг режима передачи данных.
	// - истина - данные передаются в синхронном режиме, издатель ожидает ответ.
	// - ложь   - данные передаются в асинхронном режиме.
	IsSync() bool

	// DoneSet Установка флага "обработка данных завершена".
	DoneSet()

	// Done Канал сигнала окончания обработки данных и агрегации результата.
	Done() <-chan struct{}

	// Result Возвращение результатов обработки данных.
	Result() Result

	// Context Контекст.
	Context() context.Context

	// Reset Очистка всех данных объекта.
	Reset()
}

// Result Интерфейс результата обработки данных.
type Result interface {
	// DataPut Добавление данных с результатом.
	DataPut(d ...any)

	// DataGet Возвращение данных с результатом.
	DataGet() []any

	// ErrPut Добавление ошибок в результат.
	ErrPut(e ...error)

	// ErrGet Возвращение ошибок результата.
	ErrGet(e ...error) []error
}

// Объект сущности пакета.
type impl struct {
	sync.Pool
	debug     bool       // Флаг режима отладки.
	statistic *Statistic // Статистика объектов.
}

// Обёртка над передаваемыми данными.
type data struct {
	data     any             // Данные.
	dataType reflect.Type    // Тип данных.
	sync     bool            // Флаг синхронной передачи.
	done     chan struct{}   // Канал обратной связи для синхронного вызова.
	result   *result         // Результат обработки данных.
	ctx      context.Context // Контекст.
}

// Обёртка над результатом передачи данных потребителю.
type result struct {
	data []any
	errs []error
}

// Statistic Статистика работы бассейна. Ведётся только в режиме отладки. New(isDebug=true).
type Statistic struct {
	// Создано объектов обёртки данных.
	Constructor int64

	// Уничтожено объектов обёртки данных, сборщиком мусора.
	Destructor int64

	// Получено из бассейна объектов обёртки данных.
	GetObject int64

	// Возвращено в бассейн объектов обёртки данных.
	PutObject int64
}
