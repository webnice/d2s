package server

import (
	"bytes"

	"github.com/webnice/dic"
)

// BasicAuthConfiguration Конфигурация промежуточного слоя выполняющего проверку авторизации.
// AuthFunc - Функция будет вызвана для проверки полученных данных авторизации.
// Request  - Запрос, отправляемый пользователю с просьбой авторизоваться. Можно не указывать.
// Header   - Заголовки, передаваемые с запросом с просьбой авторизоваться. Можно не указывать.
// Body     - Тело страницы, передаваемое с запросом с просьбой авторизоваться. Можно не указывать.
type BasicAuthConfiguration struct {
	AuthFunc BasicAuthFunc          // Функция выполняющая проверку простой web авторизации.
	Request  string                 // Строка запроса, отправляемая пользователю с просьбой выполнить авторизацию.
	Header   map[dic.IHeader]string // Заголовки, добавляемые в запрос аутентификации.
	Body     *bytes.Buffer          // Тело страницы запроса аутентификации.
}

// BasicAuthResponse Структура ответа на запрос с простой авторизацией.
type BasicAuthResponse struct {
	IsCorrect bool                   // Флаг подтверждения проверки доступа.
	Code      dic.IStatus            // Интерфейс объекта кода HTTP ответа.
	Header    map[dic.IHeader]string // Заголовки добавляемые в ответ.
	Body      *bytes.Buffer          // Тело ответа.
}

// BasicAuthFunc Функция передаваемая в "промежуточный слой", выполняющая проверку переданных данных
// простой авторизации.
// Функция возвращает объект BasicAuthResponse. Если свойство объекта IsCorrect равно "ложь", тогда запрос прерывается,
// и в ответ на запрос передаются данные указанные в объекте BasicAuthResponse.
// Если в объекте BasicAuthResponse свойство Body равно nil, тогда телом ответа будет код статуса ответа.
type BasicAuthFunc func(username string, password string) BasicAuthResponse
