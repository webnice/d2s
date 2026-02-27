package ambry

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
//
//goland:noinspection GoUnusedExportedFunction
func New() Interface {
	var aby = new(impl)
	aby.items = map[key][]*item{}

	return aby
}

// Set Устанавливает единственное значение для ключа key.
// Все ранее установленные значения стираются.
func (aby *impl) Set(k any, v any) {
	aby.Lock()
	aby.items[key(k)] = []*item{{Value: v}}
	aby.Unlock()
}

// Add Добавление к ключу key нового значения.
// Если у ключа key не было значения, тогда устанавливается единственное значение.
// Если у ключа key было значение, тогда новое значение добавляется в срез значений.
func (aby *impl) Add(k any, v any) {
	var ok bool

	aby.RLock()
	_, ok = aby.items[key(k)]
	aby.RUnlock()
	if ok {
		aby.Lock()
		aby.items[key(k)] = append(aby.items[key(k)], &item{Value: v})
		aby.Unlock()
	} else {
		aby.Set(k, v)
	}
}

// Has Возвращается булево значение "истина", если ключ key, имеет хотя бы одно значение.
func (aby *impl) Has(k any) (ret bool) {
	aby.RLock()
	_, ret = aby.items[key(k)]
	aby.RUnlock()
	return
}

// Get Возвращается первое значение key, если у ключа существует хотя бы одно значение.
// В противном случае, возвращается значение nil.
func (aby *impl) Get(k any) (ret any) {
	var (
		value []*item
		ok    bool
	)

	aby.RLock()
	value, ok = aby.items[key(k)]
	aby.RUnlock()
	if !ok {
		return
	}
	ret = value[0].Value

	return
}

// Del Удаление значения ключа key.
func (aby *impl) Del(k any) (ret any) {
	var (
		value []*item
		ok    bool
	)

	aby.RLock()
	value, ok = aby.items[key(k)]
	aby.RUnlock()
	if !ok {
		return
	}
	aby.Lock()
	delete(aby.items, key(k))
	aby.Unlock()
	ret = value[0].Value

	return
}

// Keys Получение всех ключей.
// Если нет ни одного ключа, возвращается пустой срез.
func (aby *impl) Keys() (ret []any) {
	var i key

	aby.RLock()
	for i = range aby.items {
		ret = append(ret, i)
	}
	aby.RUnlock()

	return
}

// GetAll Получение всех значений ключа.
// Если у ключа нет значений, возвращается пустой срез.
func (aby *impl) GetAll(k any) (ret []any) {
	var (
		value []*item
		ok    bool
		n     int
	)

	aby.RLock()
	value, ok = aby.items[key(k)]
	aby.RUnlock()
	if !ok {
		return
	}
	for n = range value {
		ret = append(ret, value[n].Value)
	}

	return
}
