package log

import (
	"container/list"
	"fmt"
	standardLog "log"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	kitModuleBus "github.com/webnice/kit/v4/module/bus"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitTypes "github.com/webnice/kit/v4/types"
)

// New Конструктор объекта интерфейса Logger.
// out - Интерфейс вывода всех сообщений.
// bus - Интерфейс шины данных.
func New(out kitTypes.SyncWriter, bus kitModuleBus.Interface) Logger {
	var log = &logger{
		wr:          out,
		bus:         bus,
		waitCounter: 0,
		subscribers: list.New(),
	}

	log.messageInitPool()

	return log
}

// Initialization Инициализация менеджера логирования.
func (log *logger) Initialization() (err error) {
	log.stdLogInterceptions()
	if err = log.bus.Subscribe(log); err != nil {
		return
	}

	return
}

// Debug Установка флага отладки.
func (log *logger) Debug(d bool) Logger { log.debug = d; return log }

// Errors Все ошибки известного состояния, которые может вернуть приложение или функция.
func (log *logger) Errors() *Error { return Errors() }

// Перехват сообщений из стандартного пакета log.
func (log *logger) stdLogInterceptions() {
	log.originalLogWriter = standardLog.Writer()
	standardLog.SetPrefix(``)
	standardLog.SetFlags(0)
	standardLog.SetOutput(log)
}

// Восстановление ранее перехваченного стандартного лога log.
func (log *logger) stdLogRestoreDefault() {
	standardLog.SetPrefix(``)
	standardLog.SetFlags(standardLog.LstdFlags)
	switch log.originalLogWriter {
	case nil:
		standardLog.SetOutput(os.Stderr)
	default:
		standardLog.SetOutput(log.originalLogWriter)
		log.originalLogWriter = nil
	}
}

// MessageGet Получение объекта сообщений лога из бассейна объектов.
func (log *logger) MessageGet() *Message { return log.messagePool.Get().(*Message) }

// MessagePut Функция для возврата объекта сообщений лога в бассейн объектов.
// Функция используется только для возврата объекта в бассейн без вывода сообщения в лог.
func (log *logger) MessagePut(msg *Message) {
	if msg == nil {
		return
	}
	msg.Reset()
	log.messagePool.Put(msg)
}

// Message Отправка log сообщения к обработчику сообщений.
func (log *logger) Message(msg *Message) {
	var (
		err      error
		errs     []error
		fatality bool
	)

	if log.doEnd {
		return
	}
	atomic.AddInt64(&log.waitCounter, 1)
	defer func() { atomic.AddInt64(&log.waitCounter, -1) }()
	// Если в шине данных отсутствует потребитель сообщений, вывод сообщения через обработчик по умолчанию.
	if !log.bus.IsSubscriber(msg) {
		log.defaultLogHandler(msg)
		return
	}
	// Отправка сообщения в шину данных.
	switch msg.Fatality {
	// Сообщение с флагом фатальности всегда отправляется синхронно, для контролируемого завершения приложения сразу
	// после разблокировки функции.
	case true:
		fatality = msg.Fatality
		if _, errs = log.bus.PublishSync(msg); len(errs) > 0 {
			log.fatalityErrorPrintStderr(errs)
		}
	default:
		// Если передача сообщения в шину данных вызывает ошибку, значит нарушена целостность приложения, не
		// зарегистрирован потребитель сообщений логирования, а это значит что надо завершать приложение.
		if err = log.bus.PublishAsync(msg); err != nil {
			fatality = true
			log.fatalityErrorPrintStderr([]error{err})
		}
	}
	if fatality {
		os.Exit(int(eApplicationFatality))
		return
	}
}

// FlushAndClose Очистка буферизированных каналов, ожидание окончания обработки всех накопленных в памяти сообщений,
// восстановление состояния стандартного лога.
func (log *logger) FlushAndClose() {
	var err error

	// Отключение приёма новых сообщений в лог.
	log.doEnd = true
	// Ожидание обработки всех принятых в обработку сообщений.
	for {
		if atomic.LoadInt64(&log.waitCounter) == 0 {
			break
		}
		runtime.Gosched()
		<-time.After(time.Millisecond * 10)
	}
	// Сброс буфера ввода/вывода.
	_ = log.wr.Sync()
	// Отписывание от шины данных.
	if err = log.bus.Unsubscribe(log); err != nil {
		log.fatalityErrorPrintStderr([]error{err})
	}
	// Восстановление настроек стандартного лога.
	log.stdLogRestoreDefault()
}

func (log *logger) fatalityErrorPrintStderr(errs []error) {
	var n int

	for n = range errs {
		_, _ = fmt.Fprintln(log.wr, errs[n].Error())
	}
	_ = log.wr.Sync()
}

// HandlerSubscribe Регистрация обработчика сообщений лога с указанием интервала уровней сообщения лога,
// которые будут передаваться подписчику.
// Min - Минимальный уровень сообщений лога, включительно.
// Max - Максимальный уровень сообщений лога, включительно.
// В режиме отладки приложения, фильтрация по уровню сообщений лога не применяется.
func (log *logger) HandlerSubscribe(handler Handler, min, max kitModuleLogLevel.Level) (err error) {
	var (
		name  string
		elm   *list.Element
		item  *subscriber
		found bool
		ok    bool
	)

	name = getFuncFullName(handler)
	for elm = log.subscribers.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*subscriber); !ok {
			continue
		}
		if name == item.Name {
			found = true
			break
		}
	}
	if found {
		err = log.Errors().HandlerAlreadySubscribed.Bind(name)
		return
	}
	item = &subscriber{
		Func: handler,
		Name: name,
		Min:  min,
		Max:  max,
	}
	log.subscribers.PushBack(item)

	return
}

// HandlerUnsubscribe Удаление обработчика сообщений лога.
func (log *logger) HandlerUnsubscribe(handler Handler) (err error) {
	var (
		name string
		del  []*list.Element
		elm  *list.Element
		item *subscriber
		ok   bool
		n    int
	)

	name = getFuncFullName(handler)
	for elm = log.subscribers.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*subscriber); !ok {
			continue
		}
		if name == item.Name {
			del = append(del, elm)
		}
	}
	if len(del) == 0 {
		err = log.Errors().HandlerSubscriptionNotFound.Bind(name)
		return
	}
	for n = range del {
		log.subscribers.Remove(del[n])
	}

	return
}
