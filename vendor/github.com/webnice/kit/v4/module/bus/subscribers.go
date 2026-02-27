package bus

import (
	"container/list"
	"reflect"
	"strings"
	"sync"
)

// Структура объекта работы с потребителями (подписчиками).
type subscribers struct {
	items      *list.List                     // Подписчики *subscriber.
	mapOfTypes map[reflect.Type][]*subscriber // Карта типов всех подписчиков.
	sync       *sync.RWMutex                  // Блокировка работы на момент изменения подписчиков и перестроения карты.
	isAny      bool                           // Если есть подписчики принимающие любые данные, значение будет "истина".
}

// Конструктор объекта *subscribers
func newSubscribers() (ret *subscribers) {
	ret = &subscribers{
		items:      list.New(),
		mapOfTypes: make(map[reflect.Type][]*subscriber),
		sync:       new(sync.RWMutex),
	}

	return
}

// Len Количество потребителей.
func (sss *subscribers) Len() int { return sss.items.Len() }

// Store Добавление потребителя (подписчика).
func (sss *subscribers) Store(s *subscriber) {
	sss.sync.Lock()
	sss.items.PushBack(s)
	sss.makeMapOfTypes()
	sss.sync.Unlock()
}

// Delete Удаление потребителя (подписчика).
func (sss *subscribers) Delete(databuserName string) (err error) {
	var (
		elm  *list.Element
		del  []*list.Element
		item *subscriber
		ok   bool
		n    int
	)

	sss.sync.Lock()
	defer sss.sync.Unlock() // Скорость разблокировки в данном месте не существенна.
	for elm = sss.items.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*subscriber); !ok {
			continue
		}
		if strings.EqualFold(item.Name, databuserName) {
			del = append(del, elm)
		}
	}
	if len(del) == 0 {
		err = Errors().DatabusSubscribeNotFound.Bind(databuserName)
		return
	}
	for n = range del {
		sss.items.Remove(del[n])
	}
	sss.makeMapOfTypes()

	return
}

// IsExistSubscriber Проверка, есть ли потребители для получения указанного типа данных.
// Вернётся "истина", если подписчики есть.
func (sss *subscribers) IsExistSubscriber(rt reflect.Type) (ret bool) {
	sss.sync.RLock()
	_, ret = sss.mapOfTypes[rt]
	sss.sync.RUnlock()
	if !ret {
		ret = sss.isAny
	}

	return
}

// GetSubscriber Получение всех потребителей для указанного типа данных, а так же потребителей,
// которые получают все типы данные без фильтрации.
func (sss *subscribers) GetSubscriber(rt reflect.Type) (ret []*subscriber) {
	var (
		elm  *list.Element
		item *subscriber
		ok   bool
	)

	sss.sync.RLock()
	ret, ok = sss.mapOfTypes[rt]
	for elm = sss.items.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*subscriber); !ok {
			continue
		}
		if len(item.Types) == 0 {
			ret = append(ret, item)
		}
	}
	sss.sync.RUnlock()

	return
}

// Пересоздание карты типов всех потребителей (подписчиков).
func (sss *subscribers) makeMapOfTypes() {
	var (
		elm  *list.Element
		item *subscriber
		ok   bool
		n    int
	)

	sss.isAny, sss.mapOfTypes = false, make(map[reflect.Type][]*subscriber)
	for elm = sss.items.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*subscriber); !ok {
			continue
		}
		if len(item.Types) == 0 {
			sss.isAny = true
			continue
		}
		for n = range item.Types {
			if _, ok = sss.mapOfTypes[item.Types[n].OriginalType]; !ok {
				sss.mapOfTypes[item.Types[n].OriginalType] = make([]*subscriber, 0, 1)
			}
			sss.mapOfTypes[item.Types[n].OriginalType] = append(sss.mapOfTypes[item.Types[n].OriginalType], item)
		}
	}
}
