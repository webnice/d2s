package d2s

import (
	"strings"
	"testing"

	d2sTypes "github.com/webnice/d2s/v2/types"
)

func strp(s string) *string { return &s }
func u64p(v uint64) *uint64 { return &v }

func TestGeneratorTags(t *testing.T) {
	d := new(impl)
	const max = 10
	for _, tc := range []struct {
		name string
		col  *d2sTypes.ColumnInfo
		want string
	}{
		{
			name: "первичный ключ",
			col: &d2sTypes.ColumnInfo{Name: "id", TypeGolang: d2sTypes.NewBase(d2sTypes.TUint64),
				TypeFull: "BIGINT(20)", IsPrimary: true, IsAutoIncrement: true, IsUnsigned: true, Precision: u64p(20)},
			want: `db:"id"         gorm:"column:id;primaryKey;autoIncrement;not null;precision:20;type:BIGINT(20) UNSIGNED"`,
		},
		{
			name: "nullable строка с размером",
			col: &d2sTypes.ColumnInfo{Name: "name", TypeGolang: d2sTypes.NewBase(d2sTypes.TString),
				TypeFull: "VARCHAR(255)", IsNullable: true, Size: u64p(255)},
			want: `db:"name"       gorm:"column:name;default:null;size:255;type:VARCHAR(255)"`,
		},
		{
			name: "decimal с точностью и масштабом",
			col: &d2sTypes.ColumnInfo{Name: "price", TypeGolang: d2sTypes.NewBase(d2sTypes.TFloat64),
				TypeFull: "DECIMAL(10,2)", Precision: u64p(10), Scale: u64p(2), Default: strp("0.00")},
			want: `db:"price"      gorm:"column:price;not null;default:0.00;precision:10;scale:2;type:DECIMAL(10,2)"`,
		},
		{
			name: "значение по умолчанию с экранированием точки с запятой",
			col: &d2sTypes.ColumnInfo{Name: "status", TypeGolang: d2sTypes.NewBase(d2sTypes.TString),
				TypeFull: "VARCHAR(10)", Default: strp("a;b")},
			want: `db:"status"     gorm:"column:status;not null;default:a\;b;type:VARCHAR(10)"`,
		},
		{
			name: "время создания",
			col: &d2sTypes.ColumnInfo{Name: "created", TypeGolang: d2sTypes.NewBase(d2sTypes.TTime),
				TypeFull: "DATETIME", Default: strp("CURRENT_TIMESTAMP")},
			want: `db:"created"    gorm:"column:created;not null;default:CURRENT_TIMESTAMP;type:DATETIME"`,
		},
		{
			name: "без знака без первичного ключа",
			col: &d2sTypes.ColumnInfo{Name: "ucnt", TypeGolang: d2sTypes.NewBase(d2sTypes.TUint32),
				TypeFull: "INT(10)", IsUnsigned: true},
			want: `db:"ucnt"       gorm:"column:ucnt;not null;type:INT(10) UNSIGNED"`,
		},
	} {
		if got := d.GeneratorTags(tc.col, max); got != tc.want {
			t.Fatalf("%s:\nполучено  %q\nожидается %q", tc.name, got, tc.want)
		}
	}
}

func TestGeneratorColumn(t *testing.T) {
	// Переопределение через командную строку имеет наивысший приоритет.
	d := new(impl)
	d.override = map[string]string{"amount": "uint64"}
	col, _ := d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "amount", Comment: "Сумма {nix.BigInt} в копейках",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TBigInt),
	}, 10)
	if col.Type != "uint64" {
		t.Fatalf("CLI приоритет: Type = %q; ожидается uint64", col.Type)
	}
	// Маркер удаляется из комментария, даже когда победила командная строка.
	if col.Comment != "Сумма в копейках" {
		t.Fatalf("CLI приоритет: Comment = %q; ожидается 'Сумма в копейках'", col.Comment)
	}

	// Переопределение через комментарий колонки.
	d = new(impl)
	col, _ = d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "ttl", Comment: "Время жизни {nix.Duration}",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TDuration),
	}, 10)
	if col.Type != "nix.Duration" || col.Comment != "Время жизни" {
		t.Fatalf("комментарий: Type = %q, Comment = %q; ожидается nix.Duration, 'Время жизни'", col.Type, col.Comment)
	}
	// Nullable колонка: базовый тип заменяется на nix-обёртку.
	col, _ = d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "ttl", Comment: "{time.Duration}", IsNullable: true,
		TypeGolang: d2sTypes.NewBase(d2sTypes.TDuration),
	}, 10)
	if col.Type != "nix.Duration" {
		t.Fatalf("nullable: Type = %q; ожидается nix.Duration", col.Type)
	}

	// Базовый маппинг без переопределений.
	d = new(impl)
	col, _ = d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "cnt", TypeGolang: d2sTypes.NewBase(d2sTypes.TInt32),
	}, 10)
	if col.Type != "int32" {
		t.Fatalf("базовый маппинг: Type = %q; ожидается int32", col.Type)
	}

	// Переводы строк в комментарии.
	col, _ = d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "cnt", Comment: "строка1\r\nстрока2",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TInt32),
	}, 10)
	if col.Comment != "строка1\\nстрока2" {
		t.Fatalf("переводы строк: Comment = %q; ожидается 'строка1\\nстрока2'", col.Comment)
	}
}

func TestGeneratorColumnUnknownType(t *testing.T) {
	d := new(impl)
	// Неизвестный тип без переопределения — ошибка.
	_, err := d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "exotic", TypeSimple: "exotic_type",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TUnknown),
	}, 10)
	if err == nil || !strings.Contains(err.Error(), "не реализовано сопоставление типа") {
		t.Fatalf("ожидается ошибка неизвестного типа, получено: %v", err)
	}
	// Неизвестный тип с переопределением через командную строку — без ошибки.
	d.override = map[string]string{"exotic": "string"}
	col, err := d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "exotic", TypeSimple: "exotic_type",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TUnknown),
	}, 10)
	if err != nil || col.Type != "string" {
		t.Fatalf("переопределение неизвестного типа: col=%#v err=%v", col, err)
	}
	// Неизвестный тип с переопределением через комментарий — без ошибки.
	d = new(impl)
	col, err = d.generatorColumn(&d2sTypes.ColumnInfo{
		Name: "exotic", TypeSimple: "exotic_type", Comment: "{[]byte}",
		TypeGolang: d2sTypes.NewBase(d2sTypes.TUnknown),
	}, 10)
	if err != nil || col.Type != "[]byte" {
		t.Fatalf("переопределение комментарием: col=%#v err=%v", col, err)
	}
}

func TestImportBlock(t *testing.T) {
	d := new(impl)
	for _, tc := range []struct {
		json, time, nix bool
		want            string
	}{
		{false, false, false, ""},
		{false, true, false, `import "time"`},
		{true, false, false, `import "encoding/json"`},
		{false, false, true, `import "github.com/webnice/nix/v2"`},
		{true, false, true, "import (\n\t\"encoding/json\"\n\n\t\"github.com/webnice/nix/v2\"\n)"},
		{true, true, true, "import (\n\t\"encoding/json\"\n\t\"time\"\n\n\t\"github.com/webnice/nix/v2\"\n)"},
	} {
		if got := d.importBlock(tc.json, tc.time, tc.nix); got != tc.want {
			t.Fatalf("importBlock(%v,%v,%v) = %q; ожидается %q", tc.json, tc.time, tc.nix, got, tc.want)
		}
	}
}

func TestGenerator(t *testing.T) {
	d := new(impl)
	inf := &d2sTypes.TableInfo{
		Database: "testdb", Table: "message", Package: "types", Struct: "Message", Comment: "Сообщение",
		Columns: []*d2sTypes.ColumnInfo{
			{Name: "id", TypeGolang: d2sTypes.NewBase(d2sTypes.TUint64), TypeFull: "BIGINT(20)",
				IsPrimary: true, IsAutoIncrement: true},
			{Name: "created_at", TypeGolang: d2sTypes.NewBase(d2sTypes.TTime), TypeFull: "DATETIME"},
			{Name: "attr", TypeGolang: d2sTypes.NewBase(d2sTypes.TRawMessage), TypeFull: "JSON"},
			{Name: "attr_null", TypeGolang: d2sTypes.NewBase(d2sTypes.TRawMessage), TypeFull: "JSON", IsNullable: true},
			{Name: "ttl", TypeGolang: d2sTypes.NewBase(d2sTypes.TDuration), TypeFull: "TIME",
				Comment: "Время {nix.Duration}"},
		},
	}
	for _, c := range inf.Columns {
		c.TypeGolang.Nullable(c.IsNullable)
	}
	buf, err := d.Generator(inf)
	if err != nil {
		t.Fatalf("Generator() ошибка: %v", err)
	}
	content := buf.String()
	for _, want := range []string{
		"package types",
		"type Message struct",
		"import (",
		"\"encoding/json\"",
		"\"time\"",
		"\"github.com/webnice/nix/v2\"",
		"ID",
		"uint64",
		"CreatedAt",
		"time.Time",
		"Attr",
		"json.RawMessage",
		"AttrNull",
		"nix.RawMessage",
		"TTL",
		"nix.Duration",
		"// Время",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("результат генерации не содержит %q:\n%s", want, content)
		}
	}
	// Маркер {...} удалён из комментария.
	if strings.Contains(content, "{nix.Duration}") {
		t.Fatalf("маркер не удалён из комментария:\n%s", content)
	}

	// Генерация без импортов.
	inf = &d2sTypes.TableInfo{
		Database: "testdb", Table: "t", Package: "types", Struct: "T", Comment: "c",
		Columns: []*d2sTypes.ColumnInfo{
			{Name: "name", TypeGolang: d2sTypes.NewBase(d2sTypes.TString), TypeFull: "VARCHAR(10)"},
		},
	}
	buf, err = d.Generator(inf)
	if err != nil {
		t.Fatalf("Generator() ошибка: %v", err)
	}
	if strings.Contains(buf.String(), "import") {
		t.Fatalf("блок импортов без необходимости:\n%s", buf.String())
	}

	// Ошибка генерации: неизвестный тип базы данных.
	inf.Columns[0].TypeGolang = d2sTypes.NewBase(d2sTypes.TUnknown)
	inf.Columns[0].TypeSimple = "exotic"
	if _, err = d.Generator(inf); err == nil ||
		!strings.Contains(err.Error(), "не реализовано сопоставление типа") {
		t.Fatalf("ожидается ошибка неизвестного типа, получено: %v", err)
	}
}
