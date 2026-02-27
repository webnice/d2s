package trace

import "sync"

// Инициализация бассейна объектов буфера данных.
// Буфер служит для минимизации выделения и освобождения памяти при выполнении вызова функций пакета путём
// переиспользования памяти выделенной для буфера ранее. Бассейн объектов служит для предотвращения очистки
// объектов сборщиком мусора, но при нехватке объектов в бассейне, создаются новые, а при избытке, лишние отдаются
// сборщику мусора.
func init() { bufferPool = new(sync.Pool); bufferPool.New = bufferNew }

// Получение объекта буфера из бассейна объектов.
func getBuffer() *buffer { return bufferPool.Get().(*buffer) }

// Возврат объекта буфера в бассейн объектов.
func putBuffer(buf *buffer) { bufferClean(buf); bufferPool.Put(buf) }

// Конструктор объектов бассейна.
func bufferNew() any {
	var buf = &buffer{
		Byte64k:     make([]byte, 0, defaultBufferLength),  // cap=64 килобайта.
		SliceString: make([]string, 0, defaultSliceLength), // cap=16 строк.
	}

	return buf
}

// Очистка свойств объекта, перед возвратом объектов в бассейн, для корректного переиспользования без старых значений.
func bufferClean(buf *buffer) {
	buf.Int = 0
	buf.UintPtr = 0
	buf.String1 = buf.String1[:0]
	buf.String2 = buf.String2[:0]
	buf.Byte64k = buf.Byte64k[:0]
	buf.SliceString = buf.SliceString[:0]
	buf.Ok = false
	buf.RuntimeFunc = nil
}
