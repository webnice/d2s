package log

import (
	"reflect"
	"runtime"
	"strings"

	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
)

// Попытка определить уровень лога по имени функции.
func findLoglevel() (ret kitModuleLogLevel.Level) {
	const (
		stackBack = 3
		fnFatal   = "log.fatal"
		fnPanic   = "log.panic"
	)
	var (
		ok      bool
		uintPtr uintptr
		fn      *runtime.Func
		fnName  string
	)

	ret = kitModuleLogLevel.Notice
	if uintPtr, _, _, ok = runtime.Caller(stackBack); ok {
		if fn = runtime.FuncForPC(uintPtr); fn != nil {
			fnName = strings.ToLower(fn.Name())
			switch {
			case strings.Contains(fnName, fnFatal):
				ret = kitModuleLogLevel.Fatal
			case strings.Contains(fnName, fnPanic):
				ret = kitModuleLogLevel.Critical
			}
		}
	}

	return
}

// Получение полного названия функции.
func getFuncFullName(obj any) (ret string) {
	var (
		rv   reflect.Value
		rt   reflect.Type
		star string
	)

	if rv = indirect(reflect.ValueOf(obj)); !rv.CanAddr() {
		ret = runtime.FuncForPC(rv.Pointer()).Name()
	} else {
		rt = indirectType(reflect.TypeOf(obj))
		if rt.Name() == "" {
			if rt.Kind() == reflect.Pointer {
				star = "*"
			}
		}
		if rt.Name() != "" {
			if rt.PkgPath() == "" {
				ret = star + rt.Name()
			} else {
				ret = star + rt.PkgPath() + "." + rt.Name()
			}
		}
	}

	return
}

func indirect(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv
}

func indirectType(rt reflect.Type) reflect.Type {
	for rt.Kind() == reflect.Ptr || rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	return rt
}
