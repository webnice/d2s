package types

import "testing"

func TestColumnNameCamelCase(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"id", "ID"},
		{"name", "Name"},
		{"created_at", "CreatedAt"},
		{"user_id", "UserID"},
		{"api_key", "APIKey"},
		{"http_server", "HTTPServer"},
		{"userID", "UserID"},
		{"CamelCase", "CamelCase"},
		{"2fa", "TwoFa"},
		{"_", "_"},
		{"__weird", "Weird"},
		{"a1_b2", "A1B2"},
		{"a1_2", "A1_2"},
		{"a__b", "AB"},
		{"a_1", "A1"},
		{"json_data", "JSONData"},
		{"ssh_key", "SSHKey"},
		{"uid", "UID"},
		{"uuid", "UUID"},
	} {
		cn := ColumnName(tc.in)
		if got := cn.CamelCase(); got != tc.want {
			t.Fatalf("CamelCase(%q) = %q; ожидается %q", tc.in, got, tc.want)
		}
	}
}

func TestColumnNameIsAbbreviation(t *testing.T) {
	cn := ColumnName("x")
	// Попадание в список аббревиатур.
	for _, abbr := range Abbreviations {
		if !cn.isAbbreviation(abbr) {
			t.Fatalf("isAbbreviation(%q) = false; ожидается true", abbr)
		}
	}
	// Промах.
	if cn.isAbbreviation("NOTABBR") {
		t.Fatal("isAbbreviation(NOTABBR) = true; ожидается false")
	}
	// Регистронезависимость.
	if !cn.isAbbreviation("id") {
		t.Fatal("isAbbreviation(id) = false; ожидается true")
	}
}
