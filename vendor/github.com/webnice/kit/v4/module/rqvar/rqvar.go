package rqvar

import (
	"net/http"
	"reflect"
)

// Get Получение объекта сущности пакета, возвращается интерфейс пакета.
//
//goland:noinspection GoUnusedExportedFunction
func Get() Interface {
	if singleton == nil {
		singleton = newObject()
	}

	return singleton
}

// Конструктор.
func newObject() (rqv *impl) { rqv = new(impl); return }

// Load Загрузка данных из запроса в объект структуры.
func (rqv *impl) Load(request *http.Request, variable any) (err error) {
	var (
		item  reflect.Value
		itpe  reflect.Type
		field reflect.Value
		flds  []reflect.StructField
		n     int
		tmp   string
		ok    bool
	)

	item = rqv.indirectValue(reflect.ValueOf(variable))
	if err = rqv.check(item); err != nil {
		return
	}
	itpe = rqv.indirectType(item.Type())
	flds = rqv.fields(itpe)
	for n = range flds {
		field = item.FieldByName(flds[n].Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if tmp, ok = flds[n].Tag.Lookup(tagContext); ok {
			if err = rqv.loadFromContext(request, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
		if tmp, ok = flds[n].Tag.Lookup(tagHeader); ok {
			if err = rqv.loadFromHeader(request, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
		if tmp, ok = flds[n].Tag.Lookup(tagCookie); ok {
			if err = rqv.loadFromCookie(request, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
		if tmp, ok = flds[n].Tag.Lookup(tagParam); ok {
			if err = rqv.loadFromUrnParam(request, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
		if tmp, ok = flds[n].Tag.Lookup(tagUrnParam); ok {
			if err = rqv.loadFromPathParam(request, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
	}
	// Функции вызываются после обработки всей структуры другими тегами.
	for n = range flds {
		field = item.FieldByName(flds[n].Name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if tmp, ok = flds[n].Tag.Lookup(tagRqFunc); ok {
			if err = rqv.loadFromRqFunc(request, item, field, tmp); err != nil {
				return
			}
			if !field.IsZero() {
				continue
			}
		}
	}

	return
}
