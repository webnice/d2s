package level

import (
	"fmt"
	"strconv"
	"strings"
)

// Константы уровня записей логирования, совпадающие по значению с принятыми в *nix системах.
const (
	// Off Логирование отключено.
	Off Level = iota - 1

	// Fatal Уровень 0: уровень предсмертных сообщений - система не стабильна, продолжение работы невозможно.
	Fatal

	// Alert Уровень 1: уровень сообщений тревоги - система нестабильна, но может частично продолжить работу.
	Alert

	// Critical Уровень 2: уровень критических ошибок - часть функционала системы работает не корректно.
	Critical

	// Error Уровень 3: уровень не критических ошибок - ошибки не прерывающие работу приложения.
	Error

	// Warning Уровень 4: уровень сообщений с предупреждениями.
	Warning

	// Notice Уровень 5: уровень информационных сообщений, требующих повышенного внимания.
	Notice

	// Info Уровень 6: сообщения информационного характера описывающие шаги выполнения алгоритмов приложения.
	Info

	// Debug Уровень 7: уровень отладочных сообщений.
	Debug

	// Trace Уровень 8: уровень максимально подробной трассировки.
	Trace
)

const tplWrongLevel = "не верное значение уровня логирования %q, ошибка: %s"

var level = map[Level]struct {
	Name  string
	Short rune
}{
	Off:      {"off", '-'},
	Fatal:    {"fatal", 'F'},
	Alert:    {"alert", 'A'},
	Critical: {"critical", 'C'},
	Error:    {"error", 'E'},
	Warning:  {"warning", 'W'},
	Notice:   {"notice", 'N'},
	Info:     {"info", 'I'},
	Debug:    {"debug", 'D'},
	Trace:    {"trace", 'T'},
}

// Level Тип уровня логирования.
type Level int8

// String Представление уровня логирования в виде строки, реализация интерфейса Stringer.
func (l Level) String() (ret string) {
	if lvs, ok := level[l]; ok {
		ret = lvs.Name
	}

	return
}

// Bytes Представление уровня логирования в виде среза байт.
func (l Level) Bytes() (ret []byte) {
	if lvs, ok := level[l]; ok {
		ret = []byte(lvs.Name)
	}

	return
}

// Short Представление уровня логирования в виде одной прописной буквы.
func (l Level) Short() (ret rune) {
	if lvs, ok := level[l]; ok {
		ret = lvs.Short
	}

	return
}

// Int Представление уровня логирования в виде числа.
func (l Level) Int() int { return int(l) }

// MarshalText Реализация интерфейса encoding.TextMarshaler.
func (l Level) MarshalText() (ret []byte, err error) { ret = []byte(l.String()); return }

// UnmarshalText Реализация интерфейса encoding.TextUnmarshaler.
func (l *Level) UnmarshalText(b []byte) (err error) {
	var (
		src  string
		eI64 error
		sI64 int64
		key  Level
	)

	src = strings.ToLower(strings.TrimSpace(string(b)))
	if sI64, eI64 = strconv.ParseInt(src, 10, 64); eI64 == nil {
		*l = ParseLevelInt64(sI64)
		return
	}
	eI64 = fmt.Errorf(tplWrongLevel, string(b), eI64)
	for key = range level {
		if src == level[key].Name || src == string(level[key].Short) {
			*l = key
			return
		}
	}
	err = eI64

	return
}

// ParseLevelInt64 Преобразование числа в тип уровня логирования.
func ParseLevelInt64(l int64) (ret Level) {
	var key Level

	ret = Off
	for key = range level {
		if int64(key) == l {
			ret = key
		}
	}

	return
}

// ParseLevelString Преобразование строки в тип уровня логирования.
func ParseLevelString(s string) (ret Level) {
	var key Level

	s, ret = strings.TrimSpace(s), Off
	for key = range level {
		lvs := level[key]
		if strings.EqualFold(s, lvs.Name) || strings.EqualFold(s, string(lvs.Short)) {
			ret = key
			break
		}
	}

	return
}
