package bus

import (
	"reflect"
	"runtime"
)

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

func makeSubscriberTypeInfo(rt reflect.Type) (ret *subscriberTypeInfo, err error) {
	var (
		star     string
		slowpoke reflect.Type
		pt       reflect.Type
	)

	ret = new(subscriberTypeInfo)
	ret.BaseType, ret.OriginalType, ret.TypeName = rt, rt, rt.String()
	slowpoke = ret.BaseType
	for {
		if pt = ret.BaseType; pt.Kind() != reflect.Pointer {
			break
		}
		if ret.BaseType = pt.Elem(); ret.BaseType == slowpoke {
			err = Errors().DatabusRecursivePointer.Bind(ret.BaseType.String())
			return
		}
		if ret.IndirectionCount%2 == 0 {
			slowpoke = slowpoke.Elem()
		}
		ret.IndirectionCount++
	}
	if rt.Name() == "" {
		if pt = rt; pt.Kind() == reflect.Pointer {
			star = "*"
			rt = pt
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			ret.TypeName = star + rt.Name()
		} else {
			ret.TypeName = star + rt.PkgPath() + "." + rt.Name()
		}
	}

	return
}

// Закрытие канала с защитой.
func safeCloseSignalChannel(c chan struct{}) {
	defer func() { _ = recover() }()
	close(c)
}
