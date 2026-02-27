package cfg

import (
	"container/list"

	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Функция создаёт канал изменения уровня работы приложения.
// Запускает поток, выполняющий изменение уровня работы приложения и вызов функций, подписанных на события изменения
// уровня работы приложения.
func runlevelChangeFuncNew() (ret chan *runLevelUp) {
	ret = make(chan *runLevelUp)
	go singleton.runlevelChangeFunc(ret)

	return
}

// Поток изменения уровня работы приложения и вызов функций подписчиков на события изменения уровня работы приложения.
func (cfg *impl) runlevelChangeFunc(ch chan *runLevelUp) {
	var (
		msg                *runLevelUp
		oldLevel, newLevel uint16
	)

	cfg.runLevelChangeChan = ch
	// Цикл чтения канала, цикл завершится с закрытием канала.
	for msg = range cfg.runLevelChangeChan {
		if newLevel = msg.newLever; newLevel == cfg.runLevel {
			msg.done <- struct{}{}
			continue
		}
		// Закрытие регистрации компонентов.
		if !cfg.comp.IsClose && newLevel > 0 {
			cfg.comp.IsClose = true
		}
		// Применение изменения уровня работы приложения.
		oldLevel, cfg.runLevel = cfg.runLevel, newLevel
		// Для режима отладки, печать в лог информации о переключении режима работы приложения в состояние завершения.
		if cfg.Debug() && cfg.runLevel != defaultRunlevelExit && newLevel == defaultRunlevelExit {
			cfg.Log().Info(tplRunlevelExit)
		}
		// Вызов функций подписчиков на изменение уровня работы приложения.
		cfg.runlevelCallSubscribers(oldLevel, newLevel)
		// Подтверждение изменения уровня работы приложения.
		msg.done <- struct{}{}
	}
}

// Вызов функций подписчиков на изменение уровня работы приложения.
func (cfg *impl) runlevelCallSubscribers(oldLevel, newLevel uint16) {
	var (
		err     error
		element *list.Element
		item    RunlevelSubscriberFn
	)

	for element = cfg.runLevelSubscribers.Front(); element != nil; element = element.Next() {
		item = element.Value.(RunlevelSubscriberFn)
		if err = cfg.runlevelSafeCallSubscriber(item, oldLevel, newLevel); err != nil {
			cfg.Gist().ErrorAppend(err)
		}
	}
}

// Безопасный вызов функции обратного вызова подписчика на изменения уровня работы приложения.
func (cfg *impl) runlevelSafeCallSubscriber(fn RunlevelSubscriberFn, oldLevel uint16, newLevel uint16) (
	err error,
) {
	// Функция защиты от паники.
	defer func() {
		if e := recover(); e != nil {
			err = cfg.Errors().RunlevelSubscriptionPanicException.Bind(e, kitModuleTrace.StackShort())
		}
	}()
	// Вызов функции подписчика.
	fn(oldLevel, newLevel)

	return
}

// Runlevel Возвращает текущее значение уровня работы приложения.
func (cfg *impl) Runlevel() uint16 { return cfg.runLevel }

// RunlevelMap Возвращает карту, описывающую план переключения уровней работы приложения.
func (cfg *impl) RunlevelMap() []uint16 { return cfg.mapRunLevel }

// RunlevelSubscribe Подписка на события изменения уровня работы приложения.
func (cfg *impl) RunlevelSubscribe(fn RunlevelSubscriberFn) (err error) {
	var (
		element      *list.Element
		item         RunlevelSubscriberFn
		found        bool
		funcFullName string
	)

	// Проверка на ошибку разработчика
	if fn == nil {
		err = cfg.Errors().RunlevelSubscribeUnsubscribeNilFunction.Bind()
		return
	}
	funcFullName = getFuncFullName(fn)
	// Поиск подписчика среди существующих подписчиков
	for element = cfg.runLevelSubscribers.Front(); element != nil; element = element.Next() {
		item = element.Value.(RunlevelSubscriberFn)
		if getFuncFullName(item) == funcFullName {
			found = true
			break
		}
	}
	// Если подписчик уже существует, возвращение ошибки
	if found {
		err = cfg.Errors().RunlevelAlreadySubscribedFunction.Bind(funcFullName)
		return
	}
	cfg.runLevelSubscribers.PushBack(fn)

	return
}

// RunlevelUnsubscribe Отписка от событий изменения уровня работы приложения.
func (cfg *impl) RunlevelUnsubscribe(fn RunlevelSubscriberFn) (err error) {
	var (
		element      *list.Element
		item         RunlevelSubscriberFn
		forDelete    []*list.Element
		funcFullName string
		n            int
	)

	// Проверка на ошибку разработчика
	if fn == nil {
		err = cfg.Errors().RunlevelSubscribeUnsubscribeNilFunction.Bind()
		return
	}
	funcFullName = getFuncFullName(fn)
	// Поиск подписчика среди существующих подписчиков
	for element = cfg.runLevelSubscribers.Front(); element != nil; element = element.Next() {
		item = element.Value.(RunlevelSubscriberFn)
		if getFuncFullName(item) == funcFullName {
			forDelete = append(forDelete, element)
		}
	}
	// Подписчик не нашёлся
	if len(forDelete) == 0 {
		err = cfg.Errors().RunlevelSubscriptionNotFound.Bind(funcFullName)
		return
	}
	// Удаление подписчика
	for n = range forDelete {
		cfg.runLevelSubscribers.Remove(forDelete[n])
	}

	return
}

// RunlevelAutoincrement Режим работы автоматического увеличения уровня работы приложения.
func (cfg *impl) RunlevelAutoincrement() (ret bool) { return !cfg.runLevelStopAutoincrement }
