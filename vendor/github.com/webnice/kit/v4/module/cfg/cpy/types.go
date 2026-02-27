package cpy

const tagName = "cpy"

var singleton = &Cpy{}

// FilterFn Тип функции фильтрации.
// Вернётся "истина", для пропуска данных.
type FilterFn func(key any, object any) (skip bool)

// Cpy Объект сущности пакета.
type Cpy struct{}
