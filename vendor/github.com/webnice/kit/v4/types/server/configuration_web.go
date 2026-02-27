package server

import "github.com/webnice/web/v3"

// WebConfiguration Структура конфигурации WEB сервера.
type WebConfiguration struct {
	// Server Конфигурация WEB сервера.
	Server web.Configuration `yaml:"Server"`

	// DocumentRoot Корень http сервера.
	// Используется в основном для не изменяемого статического контента.
	DocumentRoot string `yaml:"DocumentRoot"`

	// Pages Расположение специализированных html шаблонов для страниц сайта.
	// Код результирующих страниц генерируется на стороне сервера с использованием шаблонизатора и
	// специальных контроллеров.
	Pages string `yaml:"Pages"`
}

// WebServers Структура конфигурации группы веб серверов.
type WebServers struct {
	WebServers []WebConfiguration `yaml:"WEBServers"`
}
