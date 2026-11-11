package types

import (
	"strings"
	"testing"
)

func TestCanonicalOverrideType(t *testing.T) {
	// Все допустимые типы в каноническом написании.
	for _, name := range []string{
		"bool", "string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128",
		"[]byte", "[]string", "time.Time", "time.Duration", "json.RawMessage",
		"nix.Bool", "nix.String", "nix.Bytes", "nix.Time", "nix.Duration",
		"nix.BigInt", "nix.RawMessage", "nix.StringSlice", "nix.UUID", "nix.IP",
		"nix.Int", "nix.Int8", "nix.Int16", "nix.Int32", "nix.Int64",
		"nix.Uint", "nix.Uint8", "nix.Uint16", "nix.Uint32", "nix.Uint64",
		"nix.Uintptr", "nix.Float32", "nix.Float64", "nix.Complex64", "nix.Complex128",
	} {
		if canonical, ok := CanonicalOverrideType(name); !ok || canonical != name {
			t.Fatalf("CanonicalOverrideType(%q) = %q, %v; ожидается %q, true", name, canonical, ok, name)
		}
		// Регистронезависимость.
		if canonical, ok := CanonicalOverrideType(strings.ToUpper(name)); !ok || canonical != name {
			t.Fatalf("CanonicalOverrideType(UPPER %q) = %q, %v; ожидается %q, true", name, canonical, ok, name)
		}
	}
	// Дженерик nix.Nullable[Т].
	for _, tc := range []struct{ in, want string }{
		{"nix.Nullable[string]", "nix.Nullable[string]"},
		{"NIX.NULLABLE[STRING]", "nix.Nullable[string]"},
		{"nix.Nullable[nix.BigInt]", "nix.Nullable[nix.BigInt]"},
		{"nix.Nullable[int64]", "nix.Nullable[int64]"},
	} {
		if canonical, ok := CanonicalOverrideType(tc.in); !ok || canonical != tc.want {
			t.Fatalf("CanonicalOverrideType(%q) = %q, %v; ожидается %q, true", tc.in, canonical, ok, tc.want)
		}
	}
	// Недопустимые типы.
	for _, name := range []string{"", "  ", "uint128", "nix.Unknown", "Nullable[string]", "nix.Nullable[", "nix.Nullable[]"} {
		if _, ok := CanonicalOverrideType(name); ok {
			t.Fatalf("CanonicalOverrideType(%q) = ok; ожидается отказ", name)
		}
	}
}

func TestOverrideGoType(t *testing.T) {
	// Базовые типы: not null и nullable варианты.
	for _, tc := range []struct {
		name              string
		notNull, nullable string
	}{
		{"bool", "bool", "nix.Bool"},
		{"string", "string", "nix.String"},
		{"int", "int", "nix.Int"},
		{"int8", "int8", "nix.Int8"},
		{"int16", "int16", "nix.Int16"},
		{"int32", "int32", "nix.Int32"},
		{"int64", "int64", "nix.Int64"},
		{"uint", "uint", "nix.Uint"},
		{"uint8", "uint8", "nix.Uint8"},
		{"uint16", "uint16", "nix.Uint16"},
		{"uint32", "uint32", "nix.Uint32"},
		{"uint64", "uint64", "nix.Uint64"},
		{"uintptr", "uintptr", "nix.Uintptr"},
		{"float32", "float32", "nix.Float32"},
		{"float64", "float64", "nix.Float64"},
		{"complex64", "complex64", "nix.Complex64"},
		{"complex128", "complex128", "nix.Complex128"},
		{"[]byte", "[]byte", "nix.Bytes"},
		{"[]string", "[]string", "nix.StringSlice"},
		{"time.Time", "time.Time", "nix.Time"},
		{"time.Duration", "time.Duration", "nix.Duration"},
		{"json.RawMessage", "json.RawMessage", "nix.RawMessage"},
		{"nix.BigInt", "nix.BigInt", "nix.BigInt"},
		{"nix.UUID", "nix.UUID", "nix.UUID"},
		{"nix.IP", "nix.IP", "nix.IP"},
	} {
		if got := OverrideGoType(tc.name, false); got != tc.notNull {
			t.Fatalf("OverrideGoType(%q, false) = %q; ожидается %q", tc.name, got, tc.notNull)
		}
		if got := OverrideGoType(tc.name, true); got != tc.nullable {
			t.Fatalf("OverrideGoType(%q, true) = %q; ожидается %q", tc.name, got, tc.nullable)
		}
	}
	// Дженерик используется как есть.
	if got := OverrideGoType("nix.Nullable[string]", true); got != "nix.Nullable[string]" {
		t.Fatalf("OverrideGoType(дженерик) = %q; ожидается nix.Nullable[string]", got)
	}
}

func TestOverrideTypeNames(t *testing.T) {
	names := OverrideTypeNames()
	if len(names) == 0 {
		t.Fatal("OverrideTypeNames() вернул пустой список")
	}
	for _, name := range names[:len(names)-2] {
		if !strings.Contains(strings.Join(names, ","), name) {
			t.Fatalf("тип %q отсутствует в списке", name)
		}
	}
	if names[len(names)-1] != "nix.Nullable[T]" {
		t.Fatalf("последний элемент списка = %q; ожидается nix.Nullable[T]", names[len(names)-1])
	}
}

func TestParseOverrideList(t *testing.T) {
	// Корректный список с пробелами и пустыми парами.
	m, err := ParseOverrideList("id:uint64; payload : nix.RawMessage ; ;amount:nix.BigInt;created_at:time.Time")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	want := map[string]string{
		"id":         "uint64",
		"payload":    "nix.RawMessage",
		"amount":     "nix.BigInt",
		"created_at": "time.Time",
	}
	if len(m) != len(want) {
		t.Fatalf("длина карты = %d; ожидается %d (%#v)", len(m), len(want), m)
	}
	for k, v := range want {
		if m[k] != v {
			t.Fatalf("m[%q] = %q; ожидается %q", k, m[k], v)
		}
	}
	// Регистронезависимость имени поля.
	m, err = ParseOverrideList("ID:uint64")
	if err != nil || m["id"] != "uint64" {
		t.Fatalf("регистронезависимость: m=%#v err=%v", m, err)
	}
	// Пустой список.
	if m, err = ParseOverrideList(""); err != nil || len(m) != 0 {
		t.Fatalf("пустой список: m=%#v err=%v", m, err)
	}
	// Только разделители.
	if m, err = ParseOverrideList(" ; ;; "); err != nil || len(m) != 0 {
		t.Fatalf("только разделители: m=%#v err=%v", m, err)
	}
	// Ошибки формата.
	for _, bad := range []string{"id", "id:", ":uint64", "id:uint64;bad"} {
		if _, err = ParseOverrideList(bad); err == nil {
			t.Fatalf("ParseOverrideList(%q) = ok; ожидается ошибка формата", bad)
		} else if !strings.Contains(err.Error(), "поле:тип") {
			t.Fatalf("ошибка %q не содержит описание формата", err)
		}
	}
	// Неизвестный тип — ошибка со списком допустимых.
	_, err = ParseOverrideList("id:uint128")
	if err == nil {
		t.Fatal("ожидается ошибка для неизвестного типа")
	}
	if !strings.Contains(err.Error(), "uint128") || !strings.Contains(err.Error(), "nix.BigInt") {
		t.Fatalf("ошибка не содержит имя типа или список допустимых: %v", err)
	}
}

func TestOverrideComment(t *testing.T) {
	for _, tc := range []struct{ comment, wantType, wantCleaned string }{
		{"Время жизни {nix.Duration}", "nix.Duration", "Время жизни"},
		{"{nix.Duration} Время жизни", "nix.Duration", "Время жизни"},
		{"Время {nix.Duration} жизни", "nix.Duration", "Время жизни"},
		{"{nix.Duration}", "nix.Duration", ""},
		{"Просто комментарий", "", "Просто комментарий"},
		{"", "", ""},
		{"{примечание} текст", "", "{примечание} текст"},
		{"{примечание} текст {NIX.BigInt} конец", "nix.BigInt", "{примечание} текст конец"},
		{"{nix.Duration} и {time.Time}", "nix.Duration", "и {time.Time}"},
		{"Один {nix.Duration} два {nix.Duration} три", "nix.Duration", "Один два {nix.Duration} три"},
		{"Регистр {NIX.DURATION} типа", "nix.Duration", "Регистр типа"},
		{"Маркер без закрытия {nix.Duration", "", "Маркер без закрытия {nix.Duration"},
		{"{ nix.Duration } с пробелами", "nix.Duration", "с пробелами"},
	} {
		gotType, gotCleaned, found := OverrideComment(tc.comment)
		if gotType != tc.wantType || gotCleaned != tc.wantCleaned {
			t.Fatalf("OverrideComment(%q) = (%q, %q); ожидается (%q, %q)",
				tc.comment, gotType, gotCleaned, tc.wantType, tc.wantCleaned)
		}
		if (tc.wantType != "") != found {
			t.Fatalf("OverrideComment(%q): found = %v; ожидается %v", tc.comment, found, tc.wantType != "")
		}
	}
}
