package rqvar

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

// Загрузка значения из "печенек" HTTP запроса.
func (rqv *impl) loadFromCookie(rq *http.Request, field reflect.Value, name string) (err error) {
	var (
		coo   *http.Cookie
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
		switch coo, err = rq.Cookie(names[n]); {
		case err == nil:
		case errors.Is(err, http.ErrNoCookie):
			err = nil
			continue
		}
		if tmp = strings.TrimSpace(coo.Value); tmp == "" {
			continue
		}
		rqv.setValue(field, tmp)
		break
	}

	return
}
