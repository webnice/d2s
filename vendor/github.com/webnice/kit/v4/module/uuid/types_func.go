package uuid

// String Строковое представление названия версии.
func (vt VersionType) String() string { return vt.name }

// String Строковое представление пространства имён.
func (ns NamespaceType) String() string { return (&uuid{data: ns}).String() }
