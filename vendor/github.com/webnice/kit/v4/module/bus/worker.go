package bus

import (
	"log"
	"runtime"

	kitModulePdw "github.com/webnice/kit/v4/module/pdw"
	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Обработчик входящих объектов канала шины данных.
func (bus *impl) worker() {
	var (
		end bool
		wdi kitModulePdw.Data
		ss  []*subscriber
	)

	for {
		if end && len(bus.databus.Bus) == 0 {
			break
		}
		select {
		case <-bus.workerContext.Done():
			end = true
			continue
		case wdi = <-bus.databus.Bus:
		}
		// Отбор потребителей подписанных на тип данных.
		ss = bus.databus.Subscribers.GetSubscriber(wdi.Type())
		// Выполнение вызовов потребителей, ожидание, сбор и агрегация результатов.
		switch wdi.IsSync() {
		case true:
			bus.workerCallerSync(wdi, ss)
		default:
			bus.workerCallerAsync(wdi, ss)
		}
		runtime.Gosched()
	}
	bus.workerOnDone <- struct{}{}
}

// Асинхронный вызов.
func (bus *impl) workerCallerAsync(wdi kitModulePdw.Data, ss []*subscriber) {
	var (
		err         error
		subscribers []*subscriber
		childCount  int64
		childDone   chan struct{}
		n           int
	)

	defer func() {
		if e := recover(); e != nil {
			err = bus.Errors().DatabusPanicException.Bind(e, kitModuleTrace.StackShort())
			log.Println(err.Error())
		}
	}()
	childDone = make(chan struct{}, len(ss))
	// Делаем копию потребителей, потому что они переданы по адресу и могут исчезнуть до окончания работы функции.
	subscribers = make([]*subscriber, len(ss))
	copy(subscribers, ss)
	// Запуск безопасного вызова потребителей.
	for n = range subscribers {
		childCount++
		go func(cd chan struct{}, s *subscriber) {
			var rsp = bus.workerSafeCall(s, wdi)
			if rsp.Err != nil {
				log.Println(rsp.Err.Error())
			}
			cd <- struct{}{}
		}(childDone, subscribers[n])
	}
	// Ожидание завершения выполнения функций всех потребителей.
	for {
		if childCount == 0 {
			break
		}
		<-childDone
		childCount -= 1
	}
	close(childDone)
	// Возврат объекта в бассейн.
	bus.databus.Wrappers.Put(wdi)
}

// Синхронный вызов, с ожиданием и агрегацией результата.
func (bus *impl) workerCallerSync(wdi kitModulePdw.Data, ss []*subscriber) {
	var (
		srsp  *workerSafeCallResponse
		asrsp []*workerSafeCallResponse
		ssRsp chan *workerSafeCallResponse
		end   bool
		n     int
	)

	ssRsp, asrsp = make(chan *workerSafeCallResponse, len(ss)), make([]*workerSafeCallResponse, 0, len(ss))
	// Запуск безопасного вызова потребителей.
	for n = range ss {
		go func(s *subscriber) { ssRsp <- bus.workerSafeCall(s, wdi) }(ss[n])
	}
	// Ожидание завершения всех вызовов потребителей, либо прерывания.
	for {
		if end {
			break
		}
		select {
		case <-wdi.Context().Done():
			// Прерывание
			wdi.Result().ErrPut(wdi.Context().Err())
			end = true
		case srsp = <-ssRsp:
			// Ответ потребителя
			asrsp = append(asrsp, srsp)
		}
		// Если все потребители ответили
		if len(asrsp) == len(ss) {
			end = true
		}
	}
	close(ssRsp)
	// Подготовка результата.
	for n = range asrsp {
		if asrsp[n].Err != nil {
			wdi.Result().ErrPut(asrsp[n].Err)
			continue
		}
		if len(asrsp[n].Resp) > 0 {
			wdi.Result().DataPut(asrsp[n].Resp...)
		}
		if len(asrsp[n].Errs) > 0 {
			wdi.Result().ErrPut(asrsp[n].Errs...)
		}
	}
	wdi.DoneSet() // Сигнал готовности результата.
}

// Безопасный вызов потребителя с защитой от паники.
func (bus *impl) workerSafeCall(s *subscriber, wdi kitModulePdw.Data) (ret *workerSafeCallResponse) {
	ret = new(workerSafeCallResponse)
	defer func() {
		if e := recover(); e != nil {
			ret.Err = bus.Errors().DatabusPanicException.Bind(e, kitModuleTrace.StackShort())
		}
	}()
	ret.Resp, ret.Errs = s.Item.Consumer(wdi.IsSync(), wdi.DataGet())

	return
}
