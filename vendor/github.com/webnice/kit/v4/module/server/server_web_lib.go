package server

// Создание объекта и возвращение интерфейса InterfaceWebLib.
func newWebLib(parent *impl) *implWebLib {
	var iwl = &implWebLib{
		parent: parent,
		funcs:  new(fnHandler),
	}

	return iwl
}

// Handler Интерфейс библиотеки ВЕБ функций Handler и HandlerFunc.
func (iwl *implWebLib) Handler() InterfaceHandlerFunc { return iwl }

// Middleware Интерфейс библиотеки ВЕБ функций Middleware.
func (iwl *implWebLib) Middleware() InterfaceMiddleware { return iwl }
