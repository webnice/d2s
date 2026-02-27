package dye

/*

Помощник для вывода в терминал ANSI графики, а так же раскраски фона и текста.

*/

import "regexp"

const tplBadParam = "не верный параметр %q"

var rexHexColor = regexp.MustCompile(`(?im)^(#[0-9a-f]{6})$`)

// Объект сущности пакета.
type impl struct {
	isBackground bool            // Флаг установки следующего цвета как цвета фона.
	isBright     bool            // Флаг установки следующего цвета как яркого цвета.
	seq          []string        // Управляющая последовательность ANSI.
	profile      terminalProfile // Профайл терминала, для конвертации цветов.
}

// Объект готовой последовательности.
type sequence string
