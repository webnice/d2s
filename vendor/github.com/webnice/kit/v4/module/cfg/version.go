package cfg

// Создание объекта и возвращение интерфейса.
func newVersion(parent *impl) Version {
	var ret = &version{parent: parent}
	return ret
}

// Major returns the major version.
func (ver *version) Major() uint64 { return ver.parent.main.Version.Major() }

// Minor returns the minor version.
func (ver *version) Minor() uint64 { return ver.parent.main.Version.Minor() }

// Patch returns the patch version.
func (ver *version) Patch() uint64 { return ver.parent.main.Version.Patch() }

// Prerelease returns the prerelease version.
func (ver *version) Prerelease() string { return ver.parent.main.Version.Prerelease() }

// Metadata returns the metadata on the version.
func (ver *version) Metadata() string { return ver.parent.main.Version.Metadata() }

// String converts a Version object to a string.
func (ver *version) String() string { return ver.parent.main.Version.String() }

// MarshalJSON implements json marshaller interface.
func (ver *version) MarshalJSON() ([]byte, error) { return ver.parent.main.Version.MarshalJSON() }

// UnmarshalJSON implements json unmarshaler interface.
func (ver *version) UnmarshalJSON(b []byte) error { return ver.parent.main.Version.UnmarshalJSON(b) }
