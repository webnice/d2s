package cfg

// Version Интерфейс доступа к методам версии в семантике "Семантическое Версионирование 2.0.0"
type Version interface {
	// Major returns the major version
	Major() uint64

	// Minor returns the minor version
	Minor() uint64

	// Patch returns the patch version
	Patch() uint64

	// Prerelease returns the prerelease version
	Prerelease() string

	// Metadata returns the metadata on the version
	Metadata() string

	// String converts a Version object to a string
	String() string

	// MarshalJSON implements JSON.Marshaler interface
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements JSON.Unmarshaler interface
	UnmarshalJSON(b []byte) error
}

// Объект сущности версии, интерфейс Version
type version struct {
	parent *impl // Адрес объекта основной сущности, интерфейс Interface
}
