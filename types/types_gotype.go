package types

// GoType Карта типов.
type GoType map[bool]string

// Nullable type to string.
func (gt *GoType) Nullable(isNulable bool) (ret string) {
	ret = (*gt)[isNulable]

	return
}
