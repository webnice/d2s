package types

import "testing"

func TestBaseString(t *testing.T) {
	// Все типы typesMap в обоих вариантах.
	for _, tc := range []struct {
		simple        string
		notNull, null string
	}{
		{TBool, "bool", "nix.Bool"},
		{TBytes, "[]byte", "nix.Bytes"},
		{TFloat32, "float32", "nix.Float32"},
		{TFloat64, "float64", "nix.Float64"},
		{TInt8, "int8", "nix.Int8"},
		{TInt16, "int16", "nix.Int16"},
		{TInt32, "int32", "nix.Int32"},
		{TInt64, "int64", "nix.Int64"},
		{TString, "string", "nix.String"},
		{TStringSlice, "[]string", "nix.StringSlice"},
		{TTime, "time.Time", "nix.Time"},
		{TDuration, "time.Duration", "nix.Duration"},
		{TBigInt, "nix.BigInt", "nix.BigInt"},
		{TRawMessage, "json.RawMessage", "nix.RawMessage"},
		{TUint8, "uint8", "nix.Uint8"},
		{TUint16, "uint16", "nix.Uint16"},
		{TUint32, "uint32", "nix.Uint32"},
		{TUint64, "uint64", "nix.Uint64"},
		{TUUID, "nix.UUID", "nix.UUID"},
		{TIP, "nix.IP", "nix.IP"},
	} {
		bt := NewBase(tc.simple)
		if got := bt.String(); got != tc.notNull {
			t.Fatalf("NewBase(%q).String() = %q; ожидается %q", tc.simple, got, tc.notNull)
		}
		bt.Nullable(true)
		if got := bt.String(); got != tc.null {
			t.Fatalf("NewBase(%q) nullable.String() = %q; ожидается %q", tc.simple, got, tc.null)
		}
		bt.Nullable(false)
		if got := bt.String(); got != tc.notNull {
			t.Fatalf("NewBase(%q) notNull.String() = %q; ожидается %q", tc.simple, got, tc.notNull)
		}
	}
	// Неизвестный простой тип — пустая строка.
	bt := NewBase("unknown_simple")
	if got := bt.String(); got != "" {
		t.Fatalf("String() неизвестного типа = %q; ожидается пустая строка", got)
	}
	// TUnknown присутствует в карте типов и возвращает пустую строку.
	bt = NewBase(TUnknown)
	if got := bt.String(); got != "" {
		t.Fatalf("String() TUnknown = %q; ожидается пустая строка", got)
	}
	// GoType.Nullable напрямую.
	gt := GoType{false: "a", true: "b"}
	if got := gt.Nullable(true); got != "b" {
		t.Fatalf("GoType.Nullable(true) = %q; ожидается b", got)
	}
	if got := gt.Nullable(false); got != "a" {
		t.Fatalf("GoType.Nullable(false) = %q; ожидается a", got)
	}
}
