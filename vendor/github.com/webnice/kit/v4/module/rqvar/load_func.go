package rqvar

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Безопасный вызов функции через reflect.
func (rqv *impl) safeCall(function reflect.Value, arg ...reflect.Value) (ret []reflect.Value, err error) {
	defer func() {
		var e any
		if e = recover(); e != nil {
			err = fmt.Errorf("Вызов функции прерван паникой:\n%s\n%s", e, kitModuleTrace.StackShort())
		}
	}()
	ret = function.Call(arg)

	return
}

// Загрузка значения с использованием вызова функции структуры с указанным именем.
func (rqv *impl) loadFromRqFunc(rq *http.Request, item reflect.Value, field reflect.Value, name string) (err error) {
	var (
		function reflect.Value
		names    []string
		rsp      []reflect.Value
		n        int
	)

	if name == keySkip {
		return
	}
	names = strings.Split(name, separatorComma)
	for n = range names {
		if names[n] = strings.TrimSpace(names[n]); names[n] == "" {
			continue
		}
		if function = item.Addr().MethodByName(names[n]); !function.IsValid() {
			err = fmt.Errorf("функция %q не найдена", names[n])
			return
		}
		if rsp, err = rqv.safeCall(function, reflect.ValueOf(rq)); err != nil {
			err = fmt.Errorf("функция %q не является типом rqvar.RqFunc либо содержит ошибку.\n%s", names[n], err)
			return
		}
		if len(rsp) != 1 {
			err = fmt.Errorf(
				"вызов функции %q вернул результаты в количестве %d, ожидалось результатов: %d",
				names[n], len(rsp), 1,
			)
			return
		}
		if rsp[0].IsZero() {
			continue
		}
		if err = rqv.set(field, rsp[0]); err != nil {
			return
		}
		if !field.IsZero() {
			break
		}
	}

	return
}
