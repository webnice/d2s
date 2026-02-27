package cfg

import (
	"database/sql"
	"encoding"
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	scannerTplScanUnknownType                    = "работа со значением \"%T\" не реализована"
	scannerTplToStringNotImplemented             = "работа со значением \"%T\" не реализована, значение по умолчанию не верно"
	scannerTplToStringMarshalTextError           = "значение с типом \"%T\" не удалось преобразовать в строку, ошибка: %s"
	scannerTplKnownTypeErrTimeDuration           = "значение интервала времени %q, не распознано, ошибка: %s"
	scannerTplKnownTypeErrTimeLocation           = "значение временной зоны %q, не распознано, ошибка: %s"
	scannerTplKnownTypeErrTimeLocationNotAddress = "присвоить значение временной зоны нельзя, переменная конфигурации не является адресом"
	scannerTplToSliceValueError                  = "не удалось добавить в срез значений %s, значение %q, ошибка: %s"
	scannerTplToSliceValueSkip                   = "не удалось добавить в срез значение %q"
	scannerTplToMapCreateKey                     = "не удалось создать значение с типом %q для ключа карты из значения %q"
	scannerTplToMapCreateValue                   = "не удалось создать значение с типом %q для значения карты из значения %q"
)

type scanner struct {
	Value reflect.Value // Тип reflect.Value объекта
}

// Конструктор объекта
func makeScanner(rv reflect.Value) (ret sql.Scanner) {
	ret = &scanner{Value: rv}
	return
}

// Scan Интерфейс sql.Scanner
func (s *scanner) Scan(src any) (err error) {
	var (
		srcS     string
		nrv      reflect.Value
		nrt      reflect.Type
		ok       bool
		tmp, arr []string
		tpm      map[string]string
		n        int
	)

	// Представляет любой простой тип в виде строки.
	if srcS, err = s.itemToString(src); err != nil {
		return
	}
	// Парсинг значений для некоторых стандартизированных типов golang.
	if ok, err = s.scanKnownType(srcS); err != nil || ok {
		return
	}
	// Присвоение получателю значения string -> простой тип
	if nrv, ok, err = s.itemToKindReflectValue(s.Value.Type(), srcS); err != nil {
		return
	}
	if ok {
		s.Value.Set(nrv)
		return
	}
	switch s.Value.Kind() {
	case reflect.Slice:
		tmp = strings.Split(srcS, ",")
		for n = range tmp {
			if tmp[n] = strings.TrimSpace(tmp[n]); tmp[n] == "" {
				continue
			}
			arr = append(arr, tmp[n])
		}
		if nrv, ok, err = s.itemToSliceReflectValue(s.Value.Type().Elem(), arr); err != nil {
			return
		}
		if ok {
			s.Value.Set(reflect.AppendSlice(s.Value, nrv))
		}
	case reflect.Array:
		tmp = strings.Split(srcS, ",")
		for n = range tmp {
			if tmp[n] = strings.TrimSpace(tmp[n]); tmp[n] == "" {
				continue
			}
			arr = append(arr, tmp[n])
		}
		if nrv, ok, err = s.itemToSliceReflectValue(s.Value.Type().Elem(), arr); err != nil {
			return
		}
		for n = 0; n < s.Value.Len() && n < nrv.Len(); n++ {
			s.Value.Index(n).Set(nrv.Index(n))
		}
	case reflect.Map:
		tmp = strings.Split(srcS, ",")
		for n = range tmp {
			if tmp[n] = strings.TrimSpace(tmp[n]); tmp[n] == "" {
				continue
			}
			arr = append(arr, tmp[n])
		}
		tpm = make(map[string]string)
		for n = range arr {
			if tmp = strings.SplitN(arr[n], "=", 2); len(tmp) > 0 {
				tmp[0] = strings.TrimSpace(tmp[0])
			}
			if len(tmp) > 1 {
				tmp[1] = strings.TrimSpace(tmp[1])
			}
			tpm[tmp[0]] = tmp[1]
		}
		nrt = reflect.TypeOf(s.Value.Interface())
		if nrv, ok, err = s.itemToMapReflectValue(nrt.Key(), s.Value.Type().Elem(), tpm); err != nil {
			return
		}
		if ok {
			s.Value.Set(nrv)
		}
	default:
		err = fmt.Errorf(scannerTplScanUnknownType, s.Value.Interface())
	}

	return
}

// Представляет любой простой тип в виде строки, а так же типы, реализующие интерфейс encoding.TextMarshaler.
func (s *scanner) itemToString(item any) (ret string, err error) {
	var (
		textMarshaler encoding.TextMarshaler
		buf           []byte
		ok            bool
	)

	switch item.(type) {
	case bool:
		ret = strconv.FormatBool(item.(bool))
	case string:
		ret = item.(string)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		ret = fmt.Sprint(item)
	default:
		if textMarshaler, ok = item.(encoding.TextMarshaler); !ok {
			err = fmt.Errorf(scannerTplToStringNotImplemented, item)
			return
		}
		if buf, err = textMarshaler.MarshalText(); err != nil {
			err = fmt.Errorf(scannerTplToStringMarshalTextError, item, err)
			return
		}
		ret = string(buf)
	}

	return
}

// Парсинг значений для некоторых стандартизированных типов golang.
func (s *scanner) scanKnownType(src string) (ok bool, err error) {
	const (
		typeTimeDuration  = "time.Duration"
		typeTimeDurationA = "*time.Duration"
		typeTimeLocation  = "time.Location"
		typeTimeLocationA = "*time.Location"
	)
	var (
		svts            string
		scanner         sql.Scanner
		textUnmarshaler encoding.TextUnmarshaler
		td              time.Duration
		tl              *time.Location
	)

	// Если получатель реализует интерфейс sql.Scanner
	if scanner, ok = s.Value.Addr().Interface().(sql.Scanner); ok {
		if err = scanner.Scan(src); err != nil {
			return
		}
		return
	}
	// Если получатель реализует интерфейс encoding.TextUnmarshaler
	if textUnmarshaler, ok = s.Value.Addr().Interface().(encoding.TextUnmarshaler); ok {
		if err = textUnmarshaler.UnmarshalText([]byte(src)); err != nil {
			return
		}
		return
	}
	switch ok, svts = true, s.Value.Type().String(); svts {
	case typeTimeDuration, typeTimeDurationA:
		if td, err = time.ParseDuration(src); err != nil {
			err = fmt.Errorf(scannerTplKnownTypeErrTimeDuration, src, err)
			return
		}
		switch svts {
		case typeTimeDuration:
			s.Value.Set(reflect.ValueOf(td))
		case typeTimeDurationA:
			s.Value.Set(reflect.ValueOf(&td))
		}
	case typeTimeLocation, typeTimeLocationA:
		if tl, err = time.LoadLocation(src); err != nil {
			err = fmt.Errorf(scannerTplKnownTypeErrTimeLocation, src, err)
			return
		}
		switch svts {
		case typeTimeLocation:
			ok, err = false, fmt.Errorf(scannerTplKnownTypeErrTimeLocationNotAddress)
		case typeTimeLocationA:
			s.Value.Set(reflect.ValueOf(tl))
		}
	default:
		ok = false
	}

	return
}

// Создание reflect.Value с типом переменной назначения
func (s *scanner) itemToKindReflectValue(rvt reflect.Type, src string) (ret reflect.Value, ok bool, err error) {
	switch ret, ok = reflect.New(rvt).Elem(), true; rvt.Kind() {
	case reflect.Bool:
		var tmp bool
		if tmp, err = strconv.ParseBool(src); err == nil {
			ret.SetBool(tmp)
		}
	case reflect.String:
		ret.SetString(src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var tmp int64
		switch rvt.Kind() {
		case reflect.Int:
			tmp, err = strconv.ParseInt(src, 10, bits.UintSize)
		case reflect.Int8:
			tmp, err = strconv.ParseInt(src, 10, 8)
		case reflect.Int16:
			tmp, err = strconv.ParseInt(src, 10, 16)
		case reflect.Int32:
			tmp, err = strconv.ParseInt(src, 10, 32)
		case reflect.Int64:
			tmp, err = strconv.ParseInt(src, 10, 64)
		}
		if err == nil {
			ret.SetInt(tmp)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var tmp uint64
		switch rvt.Kind() {
		case reflect.Uint:
			tmp, err = strconv.ParseUint(src, 10, bits.UintSize)
		case reflect.Uint8:
			tmp, err = strconv.ParseUint(src, 10, 8)
		case reflect.Uint16:
			tmp, err = strconv.ParseUint(src, 10, 16)
		case reflect.Uint32:
			tmp, err = strconv.ParseUint(src, 10, 32)
		case reflect.Uint64, reflect.Uintptr:
			tmp, err = strconv.ParseUint(src, 10, 64)
		}
		if err == nil {
			ret.SetUint(tmp)
		}
	case reflect.Float32, reflect.Float64:
		var tmp float64
		switch rvt.Kind() {
		case reflect.Float32:
			tmp, err = strconv.ParseFloat(src, 32)
		case reflect.Float64:
			tmp, err = strconv.ParseFloat(src, 64)
		}
		if err == nil {
			ret.SetFloat(tmp)
		}
	case reflect.Complex64, reflect.Complex128:
		var tmp complex128
		switch rvt.Kind() {
		case reflect.Complex64:
			tmp, err = strconv.ParseComplex(src, 64)
		case reflect.Complex128:
			tmp, err = strconv.ParseComplex(src, 128)
		}
		if err == nil {
			ret.SetComplex(tmp)
		}
	default:
		ok = false
	}

	return
}

// Создание среза с типом переменной назначения и копирование значений item с конвертацией в тип назначения
func (s *scanner) itemToSliceReflectValue(rvt reflect.Type, items []string) (ret reflect.Value, ok bool, err error) {
	var (
		v     reflect.Value
		n     int
		isPtr bool
		rvta  reflect.Type
	)

	if rvt.Kind() == reflect.Ptr {
		isPtr, rvt = true, rvt.Elem()
	}

	if !isPtr {
		ret = reflect.MakeSlice(reflect.SliceOf(rvt), len(items), len(items))
	} else {
		rvta = reflect.New(rvt).Elem().Addr().Type()
		ret = reflect.MakeSlice(reflect.SliceOf(rvta), len(items), len(items))
	}
	for n = range items {
		if v, ok, err = s.itemToKindReflectValue(rvt, items[n]); err != nil {
			err = fmt.Errorf(scannerTplToSliceValueError, rvt.String(), items[n], err)
			return
		}
		if !ok {
			err = fmt.Errorf(scannerTplToSliceValueSkip, items[n])
			return
		}
		if !isPtr {
			ret.Index(n).Set(v)
		} else {
			ret.Index(n).Set(v.Addr())
		}
	}

	return
}

// Создание карты с указанным типом ключа и типом значения, наполнение карты значениями и возврат reflect.Value
func (s *scanner) itemToMapReflectValue(key reflect.Type, value reflect.Type, items map[string]string) (
	ret reflect.Value,
	ok bool,
	err error,
) {
	var (
		mapType            reflect.Type
		srcKey, srcVal     string
		srcKeyRv, srcValRv reflect.Value
	)

	mapType = reflect.MapOf(key, value)
	ret = reflect.MakeMapWithSize(mapType, len(items))
	for srcKey = range items {
		srcVal = items[srcKey]
		if srcKeyRv, ok, err = s.itemToKindReflectValue(key, srcKey); err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf(scannerTplToMapCreateKey, key.Kind().String(), srcKey)
			return
		}
		if srcValRv, ok, err = s.itemToKindReflectValue(value, srcVal); err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf(scannerTplToMapCreateValue, value.Kind().String(), srcVal)
			return
		}
		ret.SetMapIndex(srcKeyRv, srcValRv)
	}
	ok = len(items) == ret.Len()

	return
}
