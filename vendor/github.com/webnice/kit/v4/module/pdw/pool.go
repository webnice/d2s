package pdw

import (
	"context"
	"reflect"
)

// DataPut Загрузка в объект данных, флага и контекста.
func (wdo *data) DataPut(data any, isSync bool, ctx context.Context) {
	wdo.data, wdo.sync = data, isSync
	wdo.dataType = reflect.TypeOf(wdo.data)
	if ctx != nil {
		wdo.ctx = ctx
	}
}

// DataGet Возвращение оборачиваемых данных.
func (wdo *data) DataGet() any { return wdo.data }

// Type Тип передаваемых потребителю данных.
func (wdo *data) Type() reflect.Type { return wdo.dataType }

// IsSync Флаг режима передачи данных.
// - истина - данные передаются в синхронном режиме, издатель ожидает ответ;
// - ложь   - данные передаются в асинхронном режиме;
func (wdo *data) IsSync() bool { return wdo.sync }

// DoneSet Установка флага "обработка данных завершена".
func (wdo *data) DoneSet() {
	defer func() { _ = recover() }()
	wdo.done <- struct{}{}
}

// Done Канал сигнала окончания обработки данных и агрегации результата.
func (wdo *data) Done() <-chan struct{} { return wdo.done }

// Result Возвращение результатов обработки данных.
func (wdo *data) Result() Result { return wdo.result }

// Context Контекст.
func (wdo *data) Context() context.Context { return wdo.ctx }

// Reset Очистка всех данных объекта.
func (wdo *data) Reset() {
	wdo.data = nil
	wdo.sync = false
	wdo.dataType = nil
	safeCloseSignalChannel(wdo.done)
	wdo.done = make(chan struct{})
	wdo.result.data = wdo.result.data[:0]
	wdo.result.errs = wdo.result.errs[:0]
	wdo.ctx = context.Background()
}

// DataPut Добавление данных с результатом.
func (rsp *result) DataPut(d ...any) { rsp.data = append(rsp.data, d...) }

// DataGet Возвращение данных с результатом.
func (rsp *result) DataGet() []any { return rsp.data }

// ErrPut Добавление ошибок в результат.
func (rsp *result) ErrPut(e ...error) { rsp.errs = append(rsp.errs, e...) }

// ErrGet Возвращение ошибок результата.
func (rsp *result) ErrGet(e ...error) []error { return rsp.errs }
