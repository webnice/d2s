package types // import "gopkg.in/webnice/d2s.v1/d2s/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

const (
	// DrvSqlite3 sqlite3
	DrvSqlite3 = "sqlite3"

	// DrvMysql MySQL
	DrvMysql = "mysql"

	// DrvClickhouse Yandex Clickhouse
	DrvClickhouse = "clickhouse"

	// DrvRedshift Redshift
	DrvRedshift = "redshift"

	// DrvTidb Tidb
	DrvTidb = "tidb"

	// DrvPostgres Postgres SQL
	DrvPostgres = "postgres"
)

const (
	// TBool The bool type
	TBool = `bool`

	// TBytes The slice of bytes
	TBytes = `bytes`

	// TFloat64 The float64 type
	TFloat64 = `float64`

	// TInt64 The int64 type
	TInt64 = `int64`

	// TString The string type
	TString = `string`

	// TTime The time.Time type
	TTime = `time`

	// TUint64 The uint64 type
	TUint64 = `uint64`
)

var (
	// Abbreviations Main abbreviations
	Abbreviations = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", // nolint: gochecknoglobals
		"GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS",
		"RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UI",
		"UID", "UUID", "URI", "URL", "UTF8", "VM", "XML",
	}

	// NumberToWordMap Map for conversion first symbol of a structure fields from number to word
	NumberToWordMap = map[rune]string{ // nolint: gochecknoglobals
		'0': "Nil", '1': "One", '2': "Two", '3': "Three", '4': "Four",
		'5': "Five", '6': "Six", '7': "Seven", '8': "Eight", '9': "Nine",
	}

	typesMap = map[string]*GoType{ // nolint: gochecknoglobals
		TBool:    {false: "bool", true: "nul.Bool"},
		TBytes:   {false: "[]byte", true: "nul.Bytes"},
		TFloat64: {false: "float64", true: "nul.Float64"},
		TInt64:   {false: "int64", true: "nul.Int64"},
		TString:  {false: "string", true: "nul.String"},
		TTime:    {false: "time.Time", true: "nul.Time"},
		TUint64:  {false: "uint64", true: "nul.Uint64"},
	}
)
