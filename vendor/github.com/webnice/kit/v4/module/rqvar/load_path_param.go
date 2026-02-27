package rqvar

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Загрузка значения из пути URN роутинга.
func (rqv *impl) loadFromPathParam(rq *http.Request, field reflect.Value, name string) (err error) {
	var (
		names []string
		tmp   string
		n     int
	)

	if name == keySkip {
		return
	}
	names = strings.Split(name, separatorComma)
	for n = range names {
		if names[n] = strings.TrimSpace(names[n]); names[n] == "" {
			continue
		}
		tmp = chi.URLParam(rq, names[n])
		if tmp = strings.TrimSpace(tmp); tmp == "" {
			continue
		}
		rqv.setValue(field, tmp)
		break
	}

	return
}
