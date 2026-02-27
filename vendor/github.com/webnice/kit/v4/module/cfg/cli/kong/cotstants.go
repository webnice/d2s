package kong

const (
	anonymousStruct         = "<anonymous struct>"
	interpolateValueDefault = "default"
	interpolateValueEnum    = "enum"
	panicUnsupportedPath    = "unsupported Path"
	keyBeforeResolve        = "BeforeResolve"
	keyBeforeApply          = "BeforeApply"
	keyAfterApply           = "AfterApply"
	keyVersion              = "version"
	keyEnv                  = "env"
	keyValue                = "value"
	keyString               = "string"
	keyTrue                 = "true"
	keyOne                  = "1"
	keyYes                  = "yes"
	keyFalse                = "false"
	keyZero                 = "0"
	keyNo                   = "no"
	keyInt                  = "int"
	keyUint                 = "uint"
	keyFloat                = "float"
	onyOther                = "..."
	delimiterComma          = ","
	delimiterCommaSpace     = ", "
	delimiterPoint          = "."
	delimiterDash           = "-"
	delimiterDoubleDash     = "--"
	delimiterUnderscore     = "_"
	delimiterSpace          = " "
	delimiterDollar         = "$"
	cmdWithArgs             = "withargs"
	keyError                = "error"
	keyPath                 = "path"
	keyExistingFile         = "existingfile"
	keyExistingDir          = "existingdir"
	keyCounter              = "counter"
	keyDuration             = "duration"
	keyTime                 = "time"
	unexpectedEOL           = "неожиданный EOL"

	// help constants
	helpName             = "help"
	helpHelp             = "Show context-sensitive help."
	helpOrigin           = "Show context-sensitive help."
	helpShort            = rune('h')
	helpDefaultValue     = false
	defaultIndent        = 2
	defaultColumnPadding = 4
	labelCommands        = "Commands:"
	labelFlags           = "Flags:"
	labelArguments       = "Arguments:"
)

const (
	shortUsage usageOnError = iota + 1
	fullUsage
)
