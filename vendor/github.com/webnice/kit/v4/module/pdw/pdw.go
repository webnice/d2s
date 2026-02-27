package pdw

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
)

// New Конструктор бассейна с объектами для обёртки данных.
func New(isDebug bool) Interface {
	var pdw = &impl{
		Pool:      sync.Pool{},
		debug:     isDebug,
		statistic: new(Statistic),
	}

	pdw.New = pdw.constructorWrapperData

	return pdw
}

// Конструктор объектов пула.
func (pdw *impl) constructorWrapperData() any {
	var wdo = &data{
		data:     nil,
		sync:     false,
		dataType: nil,
		done:     make(chan struct{}),
		result:   new(result),
		ctx:      context.Background(),
	}

	runtime.SetFinalizer(wdo, pdw.destructorWrapperData)
	if pdw.debug {
		atomic.AddInt64(&pdw.statistic.Constructor, 1)
	}

	return wdo
}

// Деструктор объектов бассейна.
func (pdw *impl) destructorWrapperData(wdo *data) {
	safeCloseSignalChannel(wdo.done)
	if pdw.debug {
		atomic.AddInt64(&pdw.statistic.Destructor, 1)
	}
}

// Закрытие канала с защитой от паники (закрытие закрытого канала).
func safeCloseSignalChannel(c chan struct{}) {
	defer func() { _ = recover() }()
	close(c)
}

// Debug Присвоение нового значения режима отладки.
func (pdw *impl) Debug(debug bool) Interface { pdw.debug = debug; return pdw }

// Get Получение объекта из бассейна.
func (pdw *impl) Get() (ret Data) {
	var (
		wdo *data
		ok  bool
	)

	if wdo, ok = pdw.Pool.Get().(*data); ok {
		ret = wdo
	}
	if pdw.debug {
		atomic.AddInt64(&pdw.statistic.GetObject, 1)
	}

	return
}

// Put Возвращение объекта в бассейн.
func (pdw *impl) Put(wdi Data) {
	var wdo = wdi.(*data)

	wdo.Reset()
	pdw.Pool.Put(wdo)
	if pdw.debug {
		atomic.AddInt64(&pdw.statistic.PutObject, 1)
	}
}

// Statistic Статистика работы бассейна.
// Статистика ведётся только если бассейн создан с флагом отладки New(isDebug=true).
func (pdw *impl) Statistic() (ret *Statistic) {
	if !pdw.debug {
		return
	}
	// Копирование, для исключения возможности изменения извне, а так же
	// для возможности сохранения значений на момент считывания, на стороне запрашивающего.
	ret = &Statistic{
		Constructor: pdw.statistic.Constructor,
		Destructor:  pdw.statistic.Destructor,
		GetObject:   pdw.statistic.GetObject,
		PutObject:   pdw.statistic.PutObject,
	}

	return
}
