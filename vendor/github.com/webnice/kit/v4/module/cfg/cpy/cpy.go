package cpy

import (
	"database/sql"
	"fmt"
	"reflect"
	runtimeDebug "runtime/debug"
	"strings"
)

// Errors Справочник ошибок.
func (cpy *Cpy) Errors() *Error { return Errors() }

// Copy Копирование объектов.
func (cpy *Cpy) Copy(toObj any, fromObj any, selected []string, omit []string, filter FilterFn) (err error) {
	var (
		from, to, src, dst, key reflect.Value
		fromT, toT              reflect.Type
		isSlice                 bool
		i, size                 int
	)

	// Восстановление после паники.
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Catch panic: %v\nCall stack is:\n%s", e, string(runtimeDebug.Stack()))
		}
	}()
	// Значения.
	to, from = cpy.Indirect(reflect.ValueOf(toObj)), cpy.Indirect(reflect.ValueOf(fromObj))
	if isSlice, size, err = cpy.Check(to, from); err != nil {
		return
	}
	// Типы.
	toT, fromT = cpy.IndirectType(to.Type()), cpy.IndirectType(from.Type())
	// Проверка эквивалентности карты.
	if from.Kind() == reflect.Map && to.Kind() == reflect.Map && toT.String() != fromT.String() {
		err = cpy.Errors().TypeMapNotEqual.Bind()
		return
	}
	// If possible to assign, but not use filtration.
	if from.Type().AssignableTo(to.Type()) && filter != nil && from.Kind() == reflect.Map {
		// Copy map to map with filtration.
		to.Set(reflect.MakeMap(toT))
		for _, key = range from.MapKeys() {
			if filter(cpy.Indirect(key), cpy.Indirect(from.MapIndex(key)).Interface()) {
				continue
			}
			to.SetMapIndex(key, from.MapIndex(key))
		}
		return
	} else if from.Type().AssignableTo(to.Type()) && filter == nil {
		to.Set(from)
		return
	}
	for i = 0; i < size && err == nil; i++ {
		if isSlice {
			if from.Kind() == reflect.Slice {
				src = cpy.Indirect(from.Index(i))
			} else {
				src = cpy.Indirect(from)
			}
			// Фильтрация.
			if filter != nil {
				if filter(i, src.Interface()) {
					continue
				}
			}
			dst = cpy.Indirect(reflect.New(toT).Elem())
		} else {
			src = cpy.Indirect(from)
			dst = cpy.Indirect(to)
		}
		// Copy from method to field.
		err = cpy.CopyFromMethod(toT, fromT, dst, src, selected, omit, filter)
		// Copy from field to field or method.
		if err == nil {
			err = cpy.CopyFromField(toT, fromT, dst, src, selected, omit, filter)
		}
		if isSlice {
			switch {
			case dst.Addr().Type().AssignableTo(to.Type().Elem()):
				to.Set(reflect.Append(to, dst.Addr()))
			case dst.Type().AssignableTo(to.Type().Elem()):
				to.Set(reflect.Append(to, dst))
			}
		}
	}

	return
}

// IsSkip Return true for skip field.
func (cpy *Cpy) IsSkip(selected []string, omit []string, srcName string, dstName string) (ret bool) {
	var i int

	// Only selected fields.
	if len(selected) > 0 {
		ret = true
		for i = range selected {
			if selected[i] == srcName || selected[i] == dstName {
				ret = false
			}
		}
	}
	// All fields except selected.
	if len(omit) > 0 {
		ret = false
		for i = range omit {
			if omit[i] == srcName || omit[i] == dstName {
				ret = true
			}
		}
	}

	return
}

// CopyFromField Copy from field to field or method.
func (cpy *Cpy) CopyFromField(
	_ reflect.Type,
	fromT reflect.Type,
	dst reflect.Value,
	src reflect.Value,
	selected []string,
	omit []string,
	filter FilterFn,
) (err error) {
	const paramName = "name"
	var (
		fromF            reflect.Value
		field            reflect.StructField
		srcName, dstName string
	)

	// Copy from field or method to field.
	for _, field = range cpy.Fields(fromT) {
		srcName = field.Name
		if dstName = cpy.FieldReplaceName(field, paramName); dstName == "" {
			dstName = srcName
		}
		if cpy.IsSkip(selected, omit, srcName, dstName) {
			continue
		}
		fromF = src.FieldByName(srcName)
		if fromF.IsValid() {
			err = cpy.SetToFieldOrMethod(dst, dstName, fromF, srcName, selected, omit, filter)
		}
	}

	return
}

// CopyFromMethod Copy from method to field.
func (cpy *Cpy) CopyFromMethod(
	toT reflect.Type,
	_ reflect.Type,
	dst reflect.Value,
	src reflect.Value,
	selected []string,
	omit []string,
	filter FilterFn,
) (err error) {
	const paramName = "name"
	var (
		fromM            reflect.Value
		field            reflect.StructField
		srcName, dstName string
	)

	for _, field = range cpy.Fields(toT) {
		dstName = field.Name
		if srcName = cpy.FieldReplaceName(field, paramName); srcName == "" {
			srcName = dstName
		}
		if cpy.IsSkip(selected, omit, srcName, dstName) {
			continue
		}
		fromM = src.MethodByName(srcName)
		if !fromM.IsValid() && src.CanAddr() {
			fromM = src.Addr().MethodByName(srcName)
		}
		if fromM.IsValid() {
			err = cpy.SetToFieldOrMethod(dst, dstName, fromM, srcName, selected, omit, filter)
		}
	}

	return
}

// SetToFieldOrMethod Set value to field or method.
func (cpy *Cpy) SetToFieldOrMethod(
	dst reflect.Value,
	dstName string,
	from reflect.Value,
	_ string,
	selected []string,
	omit []string,
	filter FilterFn,
) (err error) {
	const paramName = "name"
	var (
		toF, toM reflect.Value
		field    reflect.StructField
		values   []reflect.Value
		name     string
		ok       bool
	)

	// By field name.
	toF = dst.FieldByName(dstName)
	// Search by tag name.
	if !toF.IsValid() {
		for _, field = range cpy.Fields(dst.Type()) {
			if name = field.Name; cpy.FieldReplaceName(field, paramName) == dstName {
				toF = dst.FieldByName(name)
				break
			}
		}
	}
	// Если field.
	if toF.IsValid() {
		// Try to can set.
		if toF.CanSet() {
			if ok, err = cpy.Set(toF, from); !ok {
				if from.Kind() == reflect.Func &&
					from.Type().NumIn() == 0 &&
					from.Type().NumOut() >= 1 {
					if values = from.Call([]reflect.Value{}); len(values) > 0 {
						if ok, err = cpy.Set(toF, values[0]); err != nil {
							return
						}
					}
				} else {
					err = cpy.Copy(toF.Addr().Interface(), from.Interface(), selected, omit, filter)
				}
			}
		}
	} else {
		// Try to set call method.
		toM = dst.MethodByName(dstName)
		if !toM.IsValid() && dst.CanAddr() {
			toM = dst.Addr().MethodByName(dstName)
		}
		if toM.IsValid() &&
			toM.Type().NumIn() == 1 &&
			from.Type().AssignableTo(toM.Type().In(0)) {
			toM.Call([]reflect.Value{from})
		}
	}

	return
}

// Indirect value get.
func (cpy *Cpy) Indirect(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv
}

// IndirectType Indirect type get.
func (cpy *Cpy) IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	return reflectType
}

// Check if input objects is correct.
func (cpy *Cpy) Check(to reflect.Value, from reflect.Value) (isSlice bool, size int, err error) {
	if !from.IsValid() {
		err = cpy.Errors().CopyFromObjectInvalid.Bind()
	}
	if !to.CanAddr() {
		err = cpy.Errors().CopyToObjectUnaddressable.Bind()
	}
	if err != nil {
		return
	}
	size = 1
	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			size = from.Len()
		}
	}

	return
}

// Set value.
func (cpy *Cpy) Set(to reflect.Value, from reflect.Value) (ok bool, err error) {
	var scanner sql.Scanner

	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
			ok = true
		}
		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
			ok = true
		} else if scanner, ok = to.Addr().Interface().(sql.Scanner); ok {
			if err, ok = scanner.Scan(from.Interface()), false; err == nil {
				ok = true
			}
		} else if from.Kind() == reflect.Ptr {
			ok, err = cpy.Set(to, from.Elem())
		}
	}

	return
}

// Fields to StructField.
func (cpy *Cpy) Fields(rt reflect.Type) (ret []reflect.StructField) {
	var (
		i int
		v reflect.StructField
	)

	if rt = cpy.IndirectType(rt); rt.Kind() == reflect.Struct {
		for i = 0; i < rt.NumField(); i++ {
			v = rt.Field(i)
			if v.Anonymous {
				ret = append(ret, cpy.Fields(v.Type)...)
				continue
			}
			ret = append(ret, v)
		}
	}

	return
}

// FieldReplaceName Get field name from tag.
func (cpy *Cpy) FieldReplaceName(field reflect.StructField, name string) (ret string) {
	var (
		tag         string
		params, tmp []string
		i           int
	)

	if tag = field.Tag.Get(tagName); tag == "" {
		return
	}
	params = strings.Split(tag, ";")
	for i = range params {
		if tmp = strings.Split(params[i], "="); len(tmp) > 1 {
			if strings.TrimSpace(tmp[0]) == name {
				ret = strings.TrimSpace(tmp[1])
			}
		}
	}

	return
}
