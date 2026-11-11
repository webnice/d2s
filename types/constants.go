package types

const (
	// DrvSqlite3 sqlite3.
	DrvSqlite3 = "sqlite3"

	// DrvMysql MySQL.
	DrvMysql = "mysql"

	// DrvClickhouse Yandex Clickhouse.
	DrvClickhouse = "clickhouse"

	// DrvRedshift Redshift.
	DrvRedshift = "redshift"

	// DrvTidb Tidb.
	DrvTidb = "tidb"

	// DrvPostgres Postgres SQL.
	DrvPostgres = "postgres"
)

const (
	// TBool Тип bool.
	TBool = `bool`

	// TBytes Срез байт.
	TBytes = `bytes`

	// TFloat32 Тип float32.
	TFloat32 = `float32`

	// TFloat64 Тип float64.
	TFloat64 = `float64`

	// TInt8 Тип int8.
	TInt8 = `int8`

	// TInt16 Тип int16.
	TInt16 = `int16`

	// TInt32 Тип int32.
	TInt32 = `int32`

	// TInt64 Тип int64.
	TInt64 = `int64`

	// TString Тип string.
	TString = `string`

	// TStringSlice Срез строк.
	TStringSlice = `strings`

	// TTime Тип time.Time.
	TTime = `time`

	// TDuration Тип time.Duration.
	TDuration = `duration`

	// TBigInt Тип nix.BigInt.
	TBigInt = `bigint`

	// TRawMessage Тип json.RawMessage.
	TRawMessage = `rawmessage`

	// TUint8 Тип uint8.
	TUint8 = `uint8`

	// TUint16 Тип uint16.
	TUint16 = `uint16`

	// TUint32 Тип uint32.
	TUint32 = `uint32`

	// TUint64 Тип uint64.
	TUint64 = `uint64`

	// TUUID Тип nix.UUID.
	TUUID = `uuid`

	// TIP Тип nix.IP.
	TIP = `ip`

	// TUnknown Неизвестный тип базы данных.
	TUnknown = `unknown`
)

var (
	// Abbreviations Основные сокращения.
	Abbreviations = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF",
		"GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS",
		"RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UI",
		"UID", "UUID", "URI", "URL", "UTF8", "VM", "XML",
	}

	// NumberToWordMap Карта для преобразования первого символа структурных полей из числа в слово.
	NumberToWordMap = map[rune]string{
		'0': "Nil", '1': "One", '2': "Two", '3': "Three", '4': "Four",
		'5': "Five", '6': "Six", '7': "Seven", '8': "Eight", '9': "Nine",
	}

	typesMap = map[string]*GoType{
		TBool:        {false: "bool", true: "nix.Bool"},
		TBytes:       {false: "[]byte", true: "nix.Bytes"},
		TFloat32:     {false: "float32", true: "nix.Float32"},
		TFloat64:     {false: "float64", true: "nix.Float64"},
		TInt8:        {false: "int8", true: "nix.Int8"},
		TInt16:       {false: "int16", true: "nix.Int16"},
		TInt32:       {false: "int32", true: "nix.Int32"},
		TInt64:       {false: "int64", true: "nix.Int64"},
		TString:      {false: "string", true: "nix.String"},
		TStringSlice: {false: "[]string", true: "nix.StringSlice"},
		TTime:        {false: "time.Time", true: "nix.Time"},
		TDuration:    {false: "time.Duration", true: "nix.Duration"},
		TBigInt:      {false: "nix.BigInt", true: "nix.BigInt"},
		TRawMessage:  {false: "json.RawMessage", true: "nix.RawMessage"},
		TUint8:       {false: "uint8", true: "nix.Uint8"},
		TUint16:      {false: "uint16", true: "nix.Uint16"},
		TUint32:      {false: "uint32", true: "nix.Uint32"},
		TUint64:      {false: "uint64", true: "nix.Uint64"},
		TUUID:        {false: "nix.UUID", true: "nix.UUID"},
		TIP:          {false: "nix.IP", true: "nix.IP"},
	}
)
