package cfg

import (
	"fmt"
	"reflect"
	runtimeDebug "runtime/debug"
	"strings"
)

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

func reflectObject(c any) (crv reflect.Value, crt reflect.Type, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = Errors().ConfigurationApplicationPanic.Bind(e, runtimeDebug.Stack())
		}
	}()
	// Проверка на nil.
	if c == nil {
		err = Errors().ConfigurationObjectIsNil.Bind()
		return
	}
	crv = indirect(reflect.ValueOf(c))
	// Проверка на не инициализированный объект равный nil.
	if !crv.IsValid() {
		err = Errors().ConfigurationObjectIsNotValid.Bind(reflect.TypeOf(c).String())
		return
	}
	crt = indirectType(crv.Type())
	// Проверка того что объект является адресом.
	if !crv.CanAddr() {
		err = Errors().ConfigurationObjectIsNotAddress.Bind(crv.Type().String())
		return
	}

	return
}

func reflectStructObject(c any) (crv reflect.Value, crt reflect.Type, err error) {
	if crv, crt, err = reflectObject(c); err != nil {
		return
	}
	// В качестве конфигураций ожидаются только структуры.
	if crt.Kind() != reflect.Struct {
		err = Errors().ConfigurationObjectIsNotStructure.Bind(reflect.TypeOf(c).String())
		return
	}

	return
}

// Удаление всех тегов, кроме перечисленных.
func reflectCleanStructTag(src reflect.StructTag, tags ...string) (ret reflect.StructTag) {
	const tpl = " %s:\"%s\""
	var (
		val string
		buf *strings.Builder
		n   int
		ok  bool
	)

	buf = &strings.Builder{}
	for n = range tags {
		if val, ok = src.Lookup(tags[n]); ok {
			buf.WriteString(fmt.Sprintf(tpl, tags[n], strings.TrimSpace(val)))
		}
	}
	ret = reflect.StructTag(strings.TrimSpace(buf.String()))

	return
}
