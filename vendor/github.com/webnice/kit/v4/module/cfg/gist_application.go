package cfg

import (
	"errors"
	"fmt"
	"os"
	runtimeDebug "runtime/debug"
	"strings"

	"github.com/webnice/dic"
)

// Завершение приложения с кодом ошибки
func (essence *gist) exitWithCode() {
	var (
		ierr dic.IError
		tmp  string
		code uint8
		n    int
	)

	for n = range essence.parent.error {
		if tmp = strings.TrimSpace(essence.parent.error[n].Error()); tmp != "" {
			_, _ = fmt.Fprintln(essence.parent.rawWriter, tmp)
		}
		if ierr = essence.Cfg().Errors().Unbind(essence.parent.error[n]); ierr != nil {
			if ierr.CodeU8().Get() > code {
				code = ierr.CodeU8().Get()
			}
		}
	}
	if code > 0 {
		os.Exit(int(code))
	}
}

// Безопасный запуск главной функции приложения
func (essence *gist) safeLaunchApplication() {
	var (
		code uint8
		err  error
	)

	defer func() {
		if e := recover(); e != nil {
			err = essence.parent.Errors().ApplicationPanicException.Bind(e, runtimeDebug.Stack())
			essence.parent.error = append(essence.parent.error, err)
		}
	}()
	// Обеспечение выполнения одной копии приложения
	// При ошибке или попытке вызова приложения из приложения, случится паника с исключением "deadlock":
	// fatal error: all goroutines are asleep - deadlock!
	essence.parent.main.FnMutex.Lock()
	defer essence.parent.main.FnMutex.Unlock()
	// Выполнение основной функции с возвратом кода ошибки
	if code, err = essence.parent.main.Fn(); code != 0 || err != nil {
		if err == nil {
			err = errors.New("")
		}
		switch code {
		case 0:
			essence.parent.error = append(
				essence.parent.error, essence.parent.Errors().ApplicationUnknownError.Bind(err),
			)
		default:
			essence.parent.error = append(
				essence.parent.error, dic.NewError(err.Error()).CodeU8().Set(code).Bind(),
			)
		}
	}

	return
}
