package rqvar

import (
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"strconv"
)

func (rqv *impl) indirectValue(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	return rv
}

func (rqv *impl) indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Pointer || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}

	return reflectType
}

// Проверка корректности переданного объекта данных.
func (rqv *impl) check(obj reflect.Value) (err error) {
	if !obj.IsValid() {
		err = fmt.Errorf("переданный объект не корректный")
		return
	}
	if !obj.CanAddr() {
		err = fmt.Errorf("передана копия объекта, требуется ссылка на объект")
		return
	}

	return
}

// Извлечение списка полей структуры данных.
func (rqv *impl) fields(rt reflect.Type) (ret []reflect.StructField) {
	var (
		i int
		v reflect.StructField
	)

	if rt = rqv.indirectType(rt); rt.Kind() == reflect.Struct {
		for i = 0; i < rt.NumField(); i++ {
			v = rt.Field(i)
			if v.Anonymous {
				ret = append(ret, rqv.fields(v.Type)...)
				continue
			}
			ret = append(ret, v)
		}
	}

	return
}

// Установка значения полю объекта с конвертацией значения в тип поля объекта структуры.
func (rqv *impl) setValue(field reflect.Value, source string) {
	var (
		err   error
		value reflect.Value
	)

	switch field.Type().String() {
	case "net.IP":
		value = rqv.indirectValue(reflect.ValueOf(net.ParseIP(source)))
		if err = rqv.set(field, value); err != nil {
			return
		}
	default:
		value = rqv.indirectValue(reflect.ValueOf(source))
		if err = rqv.set(field, value); err != nil {
			return
		}
	}
}

// Установка значения полю объекта.
func (rqv *impl) set(to reflect.Value, from reflect.Value) (err error) {
	var (
		scanner    sql.Scanner
		tmpString  string
		tmpInt64   int64
		tmpUint64  uint64
		tmpFloat64 float64
		tmpBool    bool
		ok         bool
	)

	if !from.IsValid() {
		err = fmt.Errorf("полю объекта невозможно присвоить значение")
		return
	}
	if to.Kind() == reflect.Ptr {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		to = to.Elem()
	}
	scanner, ok = to.Addr().Interface().(sql.Scanner)
	switch {
	case to.Type().String() == "net.IP":
		switch obj := from.Interface().(type) {
		case net.IP:
			to.SetBytes(obj)
		case string:
			if tmpString, ok = from.Interface().(string); ok {
				to.SetBytes(net.ParseIP(tmpString))
			}
		}
	case from.Type().ConvertibleTo(to.Type()):
		to.Set(from.Convert(to.Type()))
	case ok:
		err = scanner.Scan(from.Interface())
	case to.Kind() == reflect.Uint64,
		to.Kind() == reflect.Uint32,
		to.Kind() == reflect.Uint16,
		to.Kind() == reflect.Uint8,
		to.Kind() == reflect.Uint:
		if tmpString, ok = from.Interface().(string); ok {
			tmpUint64, _ = strconv.ParseUint(tmpString, 10, 64)
			to.SetUint(tmpUint64)
		}
	case to.Kind() == reflect.Int64,
		to.Kind() == reflect.Int32,
		to.Kind() == reflect.Int16,
		to.Kind() == reflect.Int8,
		to.Kind() == reflect.Int:
		if tmpString, ok = from.Interface().(string); ok {
			tmpInt64, _ = strconv.ParseInt(tmpString, 10, 64)
			to.SetInt(tmpInt64)
		}
	case to.Kind() == reflect.Float64,
		to.Kind() == reflect.Float32:
		if tmpString, ok = from.Interface().(string); ok {
			tmpFloat64, _ = strconv.ParseFloat(tmpString, 64)
			to.SetFloat(tmpFloat64)
		}
	case to.Kind() == reflect.Bool:
		if tmpString, ok = from.Interface().(string); ok {
			tmpBool, _ = strconv.ParseBool(tmpString)
			to.SetBool(tmpBool)
		}
	case to.Kind() == reflect.String:
		if tmpString, ok = from.Interface().(string); ok {
			to.SetString(tmpString)
		}
	case from.Kind() == reflect.Ptr:
		err = rqv.set(to, from.Elem())
	}

	return
}
