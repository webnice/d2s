package rqvar

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/webnice/dic"
)

// Загрузка значения из заголовка HTTP запроса.
func (rqv *impl) loadFromHeader(rq *http.Request, field reflect.Value, name string) (err error) {
	const prefixBearer = "Bearer" // Префикс токена в заголовке Authorization.
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
		tmp = rq.Header.Get(names[n])
		if strings.EqualFold(names[n], dic.Header().Authorization.String()) {
			tmp = strings.TrimPrefix(tmp, prefixBearer)
		}
		if tmp = strings.TrimSpace(tmp); tmp == "" {
			continue
		}
		rqv.setValue(field, tmp)
		break
	}

	return
}
