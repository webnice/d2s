package server

import "github.com/webnice/kit/v4/module/ambry"

// Создание объекта интерфейса InterfaceWebError.
func newWebError(parent *impl) *implWebError {
	var iwe = &implWebError{
		parent: parent,
		errors: ambry.New(),
	}

	return iwe
}

// Установка и получение значений ключа.
func (iwe *implWebError) do(key uint32, value error) (ret error) {
	var tmp any

	if value != nil {
		iwe.errors.Set(key, value)
	}
	if tmp = iwe.errors.Get(key); tmp != nil {
		ret, _ = tmp.(error)
	}

	return
}

// Reset Сброс всех установленных ошибок.
func (iwe *implWebError) Reset() { iwe.errors = ambry.New() }

// InternalServerError Установка и получение значения ошибки "Внутренняя ошибка сервера".
func (iwe *implWebError) InternalServerError(err error) error {
	return iwe.do(errInternalServerError, err)
}

// MethodNotAllowed Установка и получение значения ошибки "Метод не разрешён".
func (iwe *implWebError) MethodNotAllowed(err error) error {
	return iwe.do(errMethodNotAllowed, err)
}
