// Package types
package types

// GoType Golang types map
type GoType map[bool]string

// Nullable type to string
func (gt *GoType) Nullable(isNulable bool) (ret string) {
	ret = (*gt)[isNulable]
	return
}
