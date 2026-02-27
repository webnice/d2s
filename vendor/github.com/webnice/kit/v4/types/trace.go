package types

import (
	"bytes"
	"sync"
)

// TraceInfo Информации о текущем вызове и стеке.
type TraceInfo struct {
	Mux           *sync.Mutex   `json:"-"`              // Защита от одновременного доступа.
	StackTrace    *bytes.Buffer `json:"stack_trace"`    // Стек вызовов активного процесса, обрезанный до функции вызова.
	FilenameLong  *bytes.Buffer `json:"filename_long"`  // Путь и имя файла приложения из которого был совершён вызов.
	FilenameShort *bytes.Buffer `json:"filename_short"` // Название файла из которого был совершён вызов.
	Function      *bytes.Buffer `json:"function"`       // Название функции совершившей вызов.
	Package       *bytes.Buffer `json:"package"`        // Название пакета файла.
	Line          int           `json:"line"`           // Номер строки файла из которого был совершён вызов.
}

// NewTraceInfo Создание нового объекта TraceInfo.
func NewTraceInfo() (ret *TraceInfo) {
	ret = &TraceInfo{
		Mux:           new(sync.Mutex),
		StackTrace:    &bytes.Buffer{},
		FilenameLong:  &bytes.Buffer{},
		FilenameShort: &bytes.Buffer{},
		Function:      &bytes.Buffer{},
		Package:       &bytes.Buffer{},
		Line:          0,
	}

	return
}

// Copy Копирование данных передаваемого объекта.
func (ti *TraceInfo) Copy(src *TraceInfo) {
	ti.Mux.Lock()
	_, _ = ti.StackTrace.Write(src.StackTrace.Bytes())
	_, _ = ti.FilenameLong.Write(src.FilenameLong.Bytes())
	_, _ = ti.FilenameShort.Write(src.FilenameShort.Bytes())
	_, _ = ti.Function.Write(src.Function.Bytes())
	_, _ = ti.Package.Write(src.Package.Bytes())
	ti.Line = src.Line
	ti.Mux.Unlock()
}

// Reset Очистка всех полей объекта TraceInfo.
func (ti *TraceInfo) Reset() {
	ti.Mux.Lock()
	ti.StackTrace.Reset()
	ti.FilenameLong.Reset()
	ti.FilenameShort.Reset()
	ti.Function.Reset()
	ti.Package.Reset()
	ti.Line = 0
	ti.Mux.Unlock()
}
