package types // import "gopkg.in/webnice/d2s.v1/d2s/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// GoType Golang types map
type GoType map[bool]string

// Nullable type to string
func (gt *GoType) Nullable(isNulable bool) (ret string) {
	ret = (*gt)[isNulable]
	return
}
