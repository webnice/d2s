package d2s // import "gopkg.in/webnice/d2s.v1/d2s"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import ()

// New creates a new object and return interface
func New() Interface {
	var d2s = new(impl)
	return d2s
}
