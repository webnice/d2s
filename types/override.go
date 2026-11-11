package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// OverrideType Пара типов golang для переопределения: not null и nullable варианты.
type OverrideType struct {
	NotNull  string // Тип golang для not null колонки.
	Nullable string // Тип golang для nullable колонки.
}

// overrideTypes Реестр допустимых типов переопределения.
// Ключи в нижнем регистре, значение — каноническое имя и пара типов golang.
// Для базовых типов nullable вариант — соответствующая обёртка nix/v2.
// Для явных типов nix — тип используется как есть для обоих вариантов.
var overrideTypes = map[string]struct {
	canonical string
	types     OverrideType
}{
	"bool":            {"bool", OverrideType{"bool", "nix.Bool"}},
	"string":          {"string", OverrideType{"string", "nix.String"}},
	"int":             {"int", OverrideType{"int", "nix.Int"}},
	"int8":            {"int8", OverrideType{"int8", "nix.Int8"}},
	"int16":           {"int16", OverrideType{"int16", "nix.Int16"}},
	"int32":           {"int32", OverrideType{"int32", "nix.Int32"}},
	"int64":           {"int64", OverrideType{"int64", "nix.Int64"}},
	"uint":            {"uint", OverrideType{"uint", "nix.Uint"}},
	"uint8":           {"uint8", OverrideType{"uint8", "nix.Uint8"}},
	"uint16":          {"uint16", OverrideType{"uint16", "nix.Uint16"}},
	"uint32":          {"uint32", OverrideType{"uint32", "nix.Uint32"}},
	"uint64":          {"uint64", OverrideType{"uint64", "nix.Uint64"}},
	"uintptr":         {"uintptr", OverrideType{"uintptr", "nix.Uintptr"}},
	"float32":         {"float32", OverrideType{"float32", "nix.Float32"}},
	"float64":         {"float64", OverrideType{"float64", "nix.Float64"}},
	"complex64":       {"complex64", OverrideType{"complex64", "nix.Complex64"}},
	"complex128":      {"complex128", OverrideType{"complex128", "nix.Complex128"}},
	"[]byte":          {"[]byte", OverrideType{"[]byte", "nix.Bytes"}},
	"[]string":        {"[]string", OverrideType{"[]string", "nix.StringSlice"}},
	"time.time":       {"time.Time", OverrideType{"time.Time", "nix.Time"}},
	"time.duration":   {"time.Duration", OverrideType{"time.Duration", "nix.Duration"}},
	"json.rawmessage": {"json.RawMessage", OverrideType{"json.RawMessage", "nix.RawMessage"}},

	"nix.bool":        {"nix.Bool", OverrideType{"nix.Bool", "nix.Bool"}},
	"nix.string":      {"nix.String", OverrideType{"nix.String", "nix.String"}},
	"nix.bytes":       {"nix.Bytes", OverrideType{"nix.Bytes", "nix.Bytes"}},
	"nix.time":        {"nix.Time", OverrideType{"nix.Time", "nix.Time"}},
	"nix.duration":    {"nix.Duration", OverrideType{"nix.Duration", "nix.Duration"}},
	"nix.bigint":      {"nix.BigInt", OverrideType{"nix.BigInt", "nix.BigInt"}},
	"nix.rawmessage":  {"nix.RawMessage", OverrideType{"nix.RawMessage", "nix.RawMessage"}},
	"nix.stringslice": {"nix.StringSlice", OverrideType{"nix.StringSlice", "nix.StringSlice"}},
	"nix.uuid":        {"nix.UUID", OverrideType{"nix.UUID", "nix.UUID"}},
	"nix.ip":          {"nix.IP", OverrideType{"nix.IP", "nix.IP"}},
	"nix.int":         {"nix.Int", OverrideType{"nix.Int", "nix.Int"}},
	"nix.int8":        {"nix.Int8", OverrideType{"nix.Int8", "nix.Int8"}},
	"nix.int16":       {"nix.Int16", OverrideType{"nix.Int16", "nix.Int16"}},
	"nix.int32":       {"nix.Int32", OverrideType{"nix.Int32", "nix.Int32"}},
	"nix.int64":       {"nix.Int64", OverrideType{"nix.Int64", "nix.Int64"}},
	"nix.uint":        {"nix.Uint", OverrideType{"nix.Uint", "nix.Uint"}},
	"nix.uint8":       {"nix.Uint8", OverrideType{"nix.Uint8", "nix.Uint8"}},
	"nix.uint16":      {"nix.Uint16", OverrideType{"nix.Uint16", "nix.Uint16"}},
	"nix.uint32":      {"nix.Uint32", OverrideType{"nix.Uint32", "nix.Uint32"}},
	"nix.uint64":      {"nix.Uint64", OverrideType{"nix.Uint64", "nix.Uint64"}},
	"nix.uintptr":     {"nix.Uintptr", OverrideType{"nix.Uintptr", "nix.Uintptr"}},
	"nix.float32":     {"nix.Float32", OverrideType{"nix.Float32", "nix.Float32"}},
	"nix.float64":     {"nix.Float64", OverrideType{"nix.Float64", "nix.Float64"}},
	"nix.complex64":   {"nix.Complex64", OverrideType{"nix.Complex64", "nix.Complex64"}},
	"nix.complex128":  {"nix.Complex128", OverrideType{"nix.Complex128", "nix.Complex128"}},
}

// overrideNullableGeneric Шаблон дженерика nix.Nullable[Т] с базовым типом Т.
var overrideNullableGeneric = regexp.MustCompile(`(?i)^nix\.Nullable\[([^\[\]]+)\]$`)

// CanonicalOverrideType Возвращает каноническое имя типа переопределения.
// Сравнение без учёта регистра. Поддерживается дженерик nix.Nullable[Т],
// где Т — базовый скалярный тип из реестра.
func CanonicalOverrideType(name string) (canonical string, ok bool) {
	var trimmed = strings.TrimSpace(name)

	if entry, found := overrideTypes[strings.ToLower(trimmed)]; found {
		return entry.canonical, true
	}
	// Дженерик nix.Nullable[Т]: тип используется как есть, Т приводится к каноническому имени.
	if m := overrideNullableGeneric.FindStringSubmatch(trimmed); m != nil {
		if inner, found := overrideTypes[strings.ToLower(m[1])]; found {
			return "nix.Nullable[" + inner.canonical + "]", true
		}
	}

	return
}

// OverrideGoType Тип golang для канонического имени переопределения.
func OverrideGoType(canonical string, nullable bool) (ret string) {
	var (
		entry struct {
			canonical string
			types     OverrideType
		}
		found bool
	)

	if entry, found = overrideTypes[strings.ToLower(canonical)]; !found {
		// Дженерик nix.Nullable[Т] — используется как есть.
		return canonical
	}
	if nullable {
		return entry.types.Nullable
	}

	return entry.types.NotNull
}

// OverrideTypeNames Список допустимых имён типов переопределения, отсортированный.
func OverrideTypeNames() (ret []string) {
	ret = make([]string, 0, len(overrideTypes))
	for _, entry := range overrideTypes {
		ret = append(ret, entry.canonical)
	}
	sort.Strings(ret)
	ret = append(ret, "nix.Nullable[T]")

	return
}

// ParseOverrideList Разбор списка переопределений типов полей из командной строки.
// Синтаксис: "поле1:тип;поле2:тип". Имя поля — имя колонки базы данных,
// сравнение без учёта регистра. Неизвестный тип — ошибка со списком допустимых.
func ParseOverrideList(list string) (ret map[string]string, err error) {
	const (
		pairSeparator = `;`
		kvSeparator   = `:`
	)

	ret = make(map[string]string)
	for _, pair := range strings.Split(list, pairSeparator) {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, kvSeparator, 2)
		if len(kv) != 2 {
			err = fmt.Errorf("неверный формат переопределения %q, ожидается \"поле:тип\"", pair)
			return
		}
		var field, typeName = strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		if field == "" || typeName == "" {
			err = fmt.Errorf("неверный формат переопределения %q, ожидается \"поле:тип\"", pair)
			return
		}
		var canonical string
		var ok bool
		if canonical, ok = CanonicalOverrideType(typeName); !ok {
			err = fmt.Errorf("неизвестный тип данных %q в переопределении поля %q, допустимые типы: %s",
				typeName, field, strings.Join(OverrideTypeNames(), ", "))
			return
		}
		ret[strings.ToLower(field)] = canonical
	}

	return
}

// overrideCommentGroup Шаблон группы {...} в комментарии колонки.
var overrideCommentGroup = regexp.MustCompile(`\{([^{}]*)\}`)

// OverrideComment Разбор переопределения типа из комментария колонки базы данных.
// Если в комментарии в любом месте встречается группа {...}, содержимое которой
// совпадает с допустимым типом данных, — это переопределение. При нескольких
// группах берётся только первая совпавшая с типом данных.
// Возвращает каноническое имя типа и комментарий без сработавшей группы.
func OverrideComment(comment string) (canonical, cleaned string, found bool) {
	cleaned = comment
	for _, group := range overrideCommentGroup.FindAllStringSubmatchIndex(comment, -1) {
		var (
			wholeStart, wholeEnd = group[0], group[1]
			innerStart, innerEnd = group[2], group[3]
		)
		var name string
		var ok bool
		if name, ok = CanonicalOverrideType(comment[innerStart:innerEnd]); !ok {
			continue
		}
		canonical, found = name, true
		// Удаление сработавшей группы из комментария.
		cleaned = strings.Join(strings.Fields(comment[:wholeStart]+" "+comment[wholeEnd:]), " ")

		return
	}

	return
}
