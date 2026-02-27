package log

import (
	"container/list"
	"fmt"
	"runtime"

	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Consumer Интерфейс types.Databuser.
func (log *logger) Consumer(isSync bool, data any) (ret []any, errs []error) {
	var (
		childDone  chan error
		childCount int64
		message    *Message
		ok         bool
		elm        *list.Element
		item       *subscriber
		chlen      int
	)

	if message, ok = data.(*Message); !ok {
		return
	}
	if chlen = log.subscribers.Len(); chlen == 0 {
		chlen = 1
	}
	childDone = make(chan error, chlen)
	// Выбор обработчиков сообщения.
	switch log.subscribers.Len() {
	// Если обработчиков сообщений лога нет, сообщение передаётся обработчику по умолчанию.
	case 0:
		log.defaultLogHandler(message)
	// Одно сообщение, множество получателей, выдаём всем одно сообщение параллельно.
	default:
		// Передача сообщения лога, всем зарегистрированным обработчикам по очереди, с фильтрацией
		// по уровню логирования.
		for elm = log.subscribers.Front(); elm != nil; elm = elm.Next() {
			if item, ok = elm.Value.(*subscriber); !ok {
				continue
			}
			// Если уровень сообщения лога меньше минимального или больше максимального, пропускаем обработчик.
			if message.Level < item.Min || message.Level > item.Max {
				continue
			}
			childCount++
			go func(done chan error, s *subscriber, msg *Message) {
				done <- log.safeCallLogHandler(s, msg)
			}(childDone, item, message)
		}
	}
	runtime.Gosched()
	// Ожидание завершения всех процессов.
	for {
		if childCount == 0 {
			break
		}
		if err := <-childDone; err != nil {
			errs = append(errs, err)
		}
		childCount--
	}
	close(childDone)
	// Отправка объекта обратно в бассейн.
	log.MessagePut(message)

	return
}

// KnownType Интерфейс types.Databuser.
func (log *logger) KnownType() (ret []any) { ret = append(ret, new(Message)); return }

// Обработчик логов по умолчанию.
func (log *logger) defaultLogHandler(msg *Message) {
	const newLine = "\n"

	switch msg.Pattern.Len() {
	case 0:
		_, _ = fmt.Fprintln(log.wr, msg.Argument...)
	default:
		_, _ = fmt.Fprintf(log.wr, msg.Pattern.String(), msg.Argument...)
		_, _ = fmt.Fprintf(log.wr, newLine)
	}
	_ = log.wr.Sync()
}

// Безопасный запуск внешнего обработчика сообщения лога.
func (log *logger) safeCallLogHandler(item *subscriber, msg *Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = log.Errors().LogPanicException.Bind(item.Name, e, kitModuleTrace.StackShort())
		}
	}()
	item.Func(msg)

	return
}
