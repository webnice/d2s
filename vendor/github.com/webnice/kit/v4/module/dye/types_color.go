package dye

import (
	"fmt"

	"github.com/webnice/kit/v4/module/dye/colorful"
)

// Color Интерфейс цвета.
type Color interface {
	// Sequence Возвращает ANSI последовательность переключения цвета.
	Sequence(isBackground bool) string

	// String Представление последовательности переключения цвета как строки.
	String() (ret string)
}

type (
	colorNone    struct{} // Отсутствие раскраски цветом. Заглушка для терминалов не поддерживающих цвет.
	colorAnsi    int      // Цвет ANSI. Цвета (0-15), как определено стандартом ANSI.
	colorAnsi256 int      // Цвет ANSI256. Цвета (16-255), как определено стандартом ANSI 256.
	colorRgb     string   // Цвет RGB. Цвет в стандарте RGB представленный в шестнадцатеричной форме в формате #000000.
)

// Тип цвета: colorNone.

// String Представление последовательности переключения цвета как строки.
func (c colorNone) String() (ret string) { return }

// Sequence Возвращает ANSI последовательность цвета.
func (c colorNone) Sequence(isBackground bool) (ret string) { return }

// Тип цвета: colorAnsi.

// String Представление последовательности переключения цвета как строки.
func (c colorAnsi) String() string { return colorAnsiHex[c] }

// Sequence Возвращает ANSI последовательность цвета.
func (c colorAnsi) Sequence(isBackground bool) (ret string) {
	var (
		backgroundMode func(int) int
		col            int
	)

	backgroundMode = func(c int) int {
		if isBackground {
			return c + 10
		}
		return c
	}
	col = int(c)
	switch col < 8 {
	case true:
		ret = fmt.Sprintf("%d", backgroundMode(col)+30)
	default:
		ret = fmt.Sprintf("%d", backgroundMode(col-8)+90)
	}

	return
}

// Тип цвета: colorAnsi256.

// String Представление последовательности переключения цвета как строки.
func (c colorAnsi256) String() string { return colorAnsiHex[c] }

// Sequence Возвращает ANSI последовательность цвета.
func (c colorAnsi256) Sequence(isBackground bool) (ret string) {
	var prefix string

	if prefix = prefixForeground; isBackground {
		prefix = prefixBackground
	}
	ret = fmt.Sprintf("%s;5;%d", prefix, c)

	return
}

// Тип цвета: colorRgb.

// String Представление последовательности переключения цвета как строки.
func (c colorRgb) String() string { return string(c) }

// Sequence Возвращает ANSI последовательность цвета.
func (c colorRgb) Sequence(isBackground bool) (ret string) {
	var (
		err    error
		f      colorful.Color
		prefix string
	)

	if f, err = colorful.Hex(string(c)); err != nil {
		return
	}
	if prefix = prefixForeground; isBackground {
		prefix = prefixBackground
	}
	ret = fmt.Sprintf("%s;2;%d;%d;%d", prefix, uint8(f.R*255), uint8(f.G*255), uint8(f.B*255))

	return
}
