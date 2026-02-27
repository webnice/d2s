package bus

import (
	"context"
	"runtime"
	"sync/atomic"

	kitModulePdw "github.com/webnice/kit/v4/module/pdw"
)

// Создание объекта и возвращение интерфейса Essence.
func newEssence(parent *impl) Essence {
	var essence = &gist{parent: parent}
	return essence
}

// Debug Присвоение нового значения режима отладки.
func (essence *gist) Debug(debug bool) Essence {
	essence.parent.debug = debug
	essence.parent.databus.Wrappers.Debug(debug)

	return essence
}

// WorkerStart Запуск потоков обработчика шины данных.
func (essence *gist) WorkerStart(workerCount int) Essence {
	var n int

	// Если обработчики уже выполняются, выходим.
	if essence.parent.workerContext != nil || essence.parent.workerCount > 0 {
		return essence
	}
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}
	// Создание контекста.
	essence.parent.workerContext, essence.parent.workerCancel = context.WithCancel(context.Background())
	// Создание канала подтверждения завершения обработчиков.
	essence.parent.workerOnDone = make(chan struct{}, workerCount)
	// Запуск обработчиков.
	for n = 0; n < workerCount; n++ {
		atomic.AddInt64(&essence.parent.workerCount, 1)
		go essence.parent.worker()
	}
	runtime.Gosched()

	return essence
}

// WorkerStop Остановка потоков обработчика шины данных с подтверждением остановки.
// Функция блокируется до подтверждения завершения последнего потока обработчика, а сами потоки не завершаются до
// тех пор, пока в шине данных есть не обработанные данные.
func (essence *gist) WorkerStop() Essence {
	var end bool

	if essence.parent.workerContext == nil || essence.parent.workerCancel == nil {
		return essence
	}
	essence.parent.workerCancel()
	// Ожидание завершения всех обработчиков.
	for {
		if end && essence.parent.workerCount <= 0 {
			break
		}
		<-essence.parent.workerOnDone
		atomic.AddInt64(&essence.parent.workerCount, -1)
		end = true
	}
	safeCloseSignalChannel(essence.parent.workerOnDone)
	essence.parent.workerContext, essence.parent.workerCancel, essence.parent.workerCount = nil, nil, 0

	return essence
}

// Statistic Статистика работы бассейна шины данных.
// Статистика ведётся только если шина данных создана с флагом отладки New(..., isDebug=true).
// Если шина данных создана без флага отладки, статистика вернёт nil.
func (essence *gist) Statistic() *kitModulePdw.Statistic {
	return essence.parent.databus.Wrappers.Statistic()
}
