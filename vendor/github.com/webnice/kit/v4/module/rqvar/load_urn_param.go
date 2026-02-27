package rqvar

import (
	"net/http"
	"reflect"
	"strings"
)

// Загрузка значения из параметров HTTP запроса.
func (rqv *impl) loadFromUrnParam(rq *http.Request, field reflect.Value, name string) (err error) {
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
		tmp = rq.URL.Query().Get(names[n])
		if tmp = strings.TrimSpace(tmp); tmp == "" {
			continue
		}
		rqv.setValue(field, tmp)
		break
	}

	return
}
