package trace

import (
	"runtime"
	"strings"

	kitTypes "github.com/webnice/kit/v4/types"
)

// Short Получение информации о текущем вызове и коротком стеке вызовов активной горутины.
// Аргументы: stackBack - количество дополнительно пропускаемых вызовов следующих за вызовами текущего пакета.
func Short(ti *kitTypes.TraceInfo, stackBack int) {
	const isFullStack = false
	var buf *buffer

	// Сброс заполняемого объекта.
	ti.Reset()
	// Установка блокировки доступа к объекту.
	ti.Mux.Lock()
	// Получение буфера из бассейна.
	buf = getBuffer()
	// Поиск первого пакета в стеке вызовов, до первого вызова текущего пакета.
	stackBack += getStackBack(buf)
	if buf.UintPtr, buf.String1, ti.Line, buf.Ok = runtime.Caller(stackBack - 1); buf.Ok {
		_, _ = ti.FilenameLong.WriteString(buf.String1)
		// Название функции вызова.
		if buf.RuntimeFunc = runtime.FuncForPC(buf.UintPtr); buf.RuntimeFunc != nil {
			_, _ = ti.Function.WriteString(buf.RuntimeFunc.Name())
		}
		// Стек вызова.
		getStack(ti, stackBack, isFullStack, buf)
		// Коррекция названия функции.
		functionCorrect(ti)
		// Получение короткого названия файла приложения из которого был совершён вызов лога.
		if buf.SliceString = strings.Split(ti.FilenameLong.String(), packageNameSeparator); len(buf.SliceString) > 0 {
			_, _ = ti.FilenameShort.WriteString(buf.SliceString[len(buf.SliceString)-1])
		}
	}
	// Возвращение буфера в бассейн.
	putBuffer(buf)
	// Снятие блокировки доступа к объекту.
	ti.Mux.Unlock()

	return
}

// Коррекция названия функции.
func functionCorrect(ti *kitTypes.TraceInfo) {
	var buf *buffer

	buf = getBuffer()
	buf.SliceString = strings.Split(ti.Function.String(), packageNameSeparator)
	if len(buf.SliceString) > 1 {
		_, _ = ti.Package.WriteString(strings.Join(buf.SliceString[:len(buf.SliceString)-1], packageNameSeparator))
		ti.Function.Reset()
		_, _ = ti.Function.WriteString(buf.SliceString[len(buf.SliceString)-1])
	}
	buf.SliceString = strings.SplitN(ti.Function.String(), `.`, 2)
	if len(buf.SliceString) == 2 {
		if ti.Package.Len() > 0 {
			_, _ = ti.Package.WriteString(packageNameSeparator)
		}
		_, _ = ti.Package.WriteString(buf.SliceString[0])
		ti.Function.Reset()
		_, _ = ti.Function.WriteString(buf.SliceString[1])
	}
	putBuffer(buf)
}

// Поиск первого пакета в стеке вызовов, до первого вызова текущего пакета.
func getStackBack(buf *buffer) (ret int) {
	// Определение текущего пакета.
	if _, buf.String1, _, buf.Ok = runtime.Caller(0); !buf.Ok {
		return
	}
	if buf.SliceString = strings.Split(buf.String1, packageNameSeparator); len(buf.SliceString) > 0 {
		buf.String2 = strings.Join(buf.SliceString[:len(buf.SliceString)-1], packageNameSeparator)
	}
	for {
		_, buf.String1, _, buf.Ok = runtime.Caller(ret)
		if strings.Contains(buf.String1, buf.String2) && buf.Ok {
			ret++
			continue
		}
		break
	}

	return
}

// GetStack Загрузка стека вызова.
func getStack(ti *kitTypes.TraceInfo, stackBack int, isFullStack bool, buf *buffer) {
	var end int

	if len(buf.Byte64k) != cap(buf.Byte64k) {
		// Увеличиваем размер до размера выделенной памяти, в буфере будут видны старые данные, но это не важно, при
		// успешном выполнении, данные будут затёрты и обрезаны.
		buf.Byte64k = buf.Byte64k[0:cap(buf.Byte64k)]
	}
	// Попытка получить стек до ошибки либо до тех пор, пока буфер будет достаточного размера для успешной загрузки.
	for {
		if buf.Int = runtime.Stack(buf.Byte64k, isFullStack); buf.Int < len(buf.Byte64k) {
			buf.Byte64k = buf.Byte64k[:buf.Int]
			break
		}
		// Нужен новый буфер в два раза больше, создание нового буфера.
		buf.Byte64k = make([]byte, 2*len(buf.Byte64k))
	}
	// Результат
	ti.StackTrace.Write(buf.Byte64k[:])
	buf.Byte64k = buf.Byte64k[:0]
	// Если не требуется получение полного стека, обрезание первых строк с не нужными пакетами.
	if !isFullStack {
		buf.SliceString = strings.Split(ti.StackTrace.String(), stackLineSeparator)
		ti.StackTrace.Reset()
		if end = 1 + (stackBack * 2); end < 0 {
			end = 0
		}
		_, _ = ti.StackTrace.WriteString(strings.Join(
			append(buf.SliceString[0:1], buf.SliceString[end:]...),
			stackLineSeparator,
		))
	}
}
