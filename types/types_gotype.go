package types

// GoType Карта типов.
type GoType map[bool]string

// Nullable Тип в строку.
func (gt *GoType) Nullable(isNulable bool) (ret string) {
	ret = (*gt)[isNulable]

	return
}
