package rqvar

import (
	"net/http"
	"reflect"
	"strings"
)

// Загрузка значения из контекста HTTP запроса.
func (rqv *impl) loadFromContext(rq *http.Request, field reflect.Value, name string) (err error) {
	var (
		names []string
		tmp   reflect.Value
		n     int
		value any
	)

	if name == keySkip {
		return
	}
	names = strings.Split(name, separatorComma)
	for n = range names {
		if names[n] = strings.TrimSpace(names[n]); names[n] == "" {
			continue
		}
		if value = rq.Context().Value(names[n]); value == nil {
			continue
		}
		tmp = reflect.ValueOf(value)
		if err = rqv.set(field, tmp); err != nil {
			err = nil
			continue
		}
		break
	}

	return
}
