package server

import (
	"fmt"

	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Вызов функции с защитой от паники.
func (iwl *implWebLib) safeCall(fn func()) (err error) {
	const errPanic = "Паника: %q\nСтек вызовов, в момент паники:\n%s."
	// Функция защиты от паники.
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(errPanic, e, kitModuleTrace.StackShort())
		}
	}()
	fn()

	return
}
