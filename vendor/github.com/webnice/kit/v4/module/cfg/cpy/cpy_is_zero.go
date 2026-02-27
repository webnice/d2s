package cpy

import (
	"database/sql"
	"fmt"
	"reflect"
	runtimeDebug "runtime/debug"
)

// CopyToIsZero Копирование значений только если значение в объекте назначения пустое.
func (cpy *Cpy) CopyToIsZero(toRv reflect.Value, fromRv reflect.Value) (err error) {
	var (
		to        reflect.Value
		from      reflect.Value
		fromT     reflect.Type
		fromV     reflect.Value
		field     reflect.StructField
		fieldName string
	)

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("catch panic: %v\ncall stack is:\n%s", e, string(runtimeDebug.Stack()))
		}
	}()
	to, from, fromT = cpy.Indirect(toRv), cpy.Indirect(fromRv), cpy.IndirectType(fromRv.Type())
	for _, field = range cpy.Fields(fromT) {
		fieldName = field.Name
		if fromV = from.FieldByName(fieldName); fromV.IsValid() && !fromV.IsZero() {
			err = cpy.SetToFieldIsZero(to, fieldName, fromV)
		}

	}

	return
}

// SetToFieldIsZero Копирование одного поля по имени, только если значение в объекте назначения пустое.
func (cpy *Cpy) SetToFieldIsZero(toRv reflect.Value, toName string, fromV reflect.Value) (err error) {
	var (
		toV, toM reflect.Value
		values   []reflect.Value
		ok       bool
	)

	switch toV = toRv.FieldByName(toName); toV.IsValid() {
	case true:
		if toV.CanSet() {
			if ok, err = cpy.SetToIsZero(toV, fromV); !ok {
				if fromV.Kind() == reflect.Func && fromV.Type().NumIn() == 0 && fromV.Type().NumOut() >= 1 {
					if values = fromV.Call([]reflect.Value{}); len(values) > 0 {
						if ok, err = cpy.SetToIsZero(toV, values[0]); err != nil {
							return
						}
					}
				} else {
					err = cpy.CopyToIsZero(toV.Addr(), fromV)
				}
			}
		}
	default:
		if toM = toRv.MethodByName(toName); !toM.IsValid() && toRv.CanAddr() {
			toM = toRv.Addr().MethodByName(toName)
		}
		if toM.IsValid() && toM.Type().NumIn() == 1 && fromV.Type().AssignableTo(toM.Type().In(0)) {
			toM.Call([]reflect.Value{fromV})
		}
	}

	return
}

// SetToIsZero Установка значения только если значение получателя является пустым значением.
func (cpy *Cpy) SetToIsZero(to reflect.Value, from reflect.Value) (ok bool, err error) {
	var scanner sql.Scanner

	if !from.IsValid() {
		return
	}
	if to.Kind() == reflect.Ptr {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		to, ok = to.Elem(), true
	}
	// Пропускаем поля с установленными не пустыми значениями
	if !to.IsZero() {
		return
	}
	switch {
	case from.Type().ConvertibleTo(to.Type()):
		to.Set(from.Convert(to.Type()))
		ok = true
	case from.Kind() == reflect.Ptr:
		ok, err = cpy.Set(to, from.Elem())
	default:
		if scanner, ok = to.Addr().Interface().(sql.Scanner); ok {
			if err, ok = scanner.Scan(from.Interface()), false; err == nil {
				ok = true
			}
		}
	}

	return
}
