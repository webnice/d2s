package trace

import (
	"runtime"
	"sync"
)

const (
	packageNameSeparator string = "/"     // Разделитель названия пакета.
	stackLineSeparator   string = "\n"    // Разделитель строк стека.
	defaultBufferLength  int    = 1 << 16 // Размер буфера по умолчанию для среза байт - cap=64 килобайта.
	defaultSliceLength   int    = 16      // Размер буфера по умолчанию для среза строк - cap=16.
)

var bufferPool *sync.Pool

// Структура объектов бассейна
type buffer struct {
	Int         int           // Переиспользуемый int.
	UintPtr     uintptr       // Переиспользуемый uintptr.
	String1     string        // Переиспользуемая строка.
	String2     string        // Переиспользуемая строка.
	Byte64k     []byte        // Переиспользуемый срез байт размером 64 килобайта.
	SliceString []string      // Переиспользуемый срез строк.
	Ok          bool          // Переиспользуемый булев.
	RuntimeFunc *runtime.Func // Переиспользуемый *runtime.Func.
}
